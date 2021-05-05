package backup

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Simonpoon93/etcd-operator/pkg/file"

	"github.com/coreos/etcd/clientv3/snapshot"

	"github.com/go-logr/logr"

	"github.com/coreos/etcd/clientv3"

	"github.com/go-logr/zapr"

	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func logErr(log logr.Logger, err error, message string) error {
	log.Error(err, message)
	return fmt.Errorf("%s: %s", message, err)
}

func main() {
	var (
		backupTempDir      string
		etcdEndpoints      string
		dialTimeoutSeconds int64
		timeoutSeconds     int64
	)
	// 定义参数
	flag.StringVar(&backupTempDir, "backup-tmp-dir", os.TempDir(), "The directory to temp place backup etcd cluster")
	flag.StringVar(&etcdEndpoints, "etcdEndpoints", "", "Etcd Endpoints")
	flag.Int64Var(&dialTimeoutSeconds, "dial-timeout-seconds", 5, "Timeout for dialing the Etcd")
	flag.Int64Var(&timeoutSeconds, "timeout-seconds", 60, "Timeout for backup the Etcd")

	// 生成带有超时的ctx
	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeoutSeconds))

	defer ctxCancel()

	log := ctrl.Log.WithName("backup")
	// 定义一个本地的数据目录
	localPath := filepath.Join(backupTempDir, "snapshot.db")

	zapLogger := zap.NewRaw(zap.UseDevMode(true))
	ctrl.SetLogger(zapr.NewLogger(zapLogger))

	// 创建etcd snapshot Client
	etcdManager := snapshot.NewV3(zapLogger)

	log.Info("Connecting to Etcd and getting snapshot data")

	// 保存etcd snapshot数据到localPath
	err := etcdManager.Save(ctx, clientv3.Config{
		Endpoints:   []string{etcdEndpoints},
		DialTimeout: time.Second * time.Duration(timeoutSeconds),
	}, localPath)
	if err != nil {
		panic(logErr(log, err, "failed to get etcd snapshot data"))
	}

	// 数据保存到本地成功 上传到其他存储中
	// TODO 根据传递进来的参数判断初始化s3还是oss
	endpoint := "play.min.io"
	accessKeyID := "Q3AM3UQ867SPQQA43P2F"
	secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
	s3Uploader := file.NewS3Uploader(endpoint, accessKeyID, secretAccessKey)
	log.Info("Uploading snapshot")

	// 上传文件到minio
	size, err := s3Uploader.Upload(ctx, localPath)
	if err != nil {
		panic(logErr(log, err, "failed to upload backup file to S3"))
	}

	log.WithValues("upload-size", size).Info("Backup completed")
}
