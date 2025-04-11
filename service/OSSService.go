package service

import (
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"note_app_server_mq/global"
)

// DeleteDir 删除文件夹
// @params bucket 桶名称
// @param dirName 目录名
func DeleteDir(ctx context.Context, bucketName, dirName string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		client := global.OssClientPool.Get().(*oss.Client)
		defer global.OssClientPool.Put(client)
		bucket, err := client.Bucket(bucketName)
		if err != nil {
			return err
		}

		var totalDeleted int
		marker := oss.Marker("")
		prefix := oss.Prefix(dirName + "/")
		for {
			lor, err1 := bucket.ListObjects(marker, prefix)
			if err1 != nil {
				log.Fatal(err1)
			}

			objects := make([]string, len(lor.Objects))
			for i, o := range lor.Objects {
				objects[i] = o.Key
			}

			delRes, err2 := bucket.DeleteObjects(objects, oss.DeleteObjectsQuiet(true))
			if err2 != nil {
				log.Fatalf("Failed to delete objects: %v", err)
			}

			if len(delRes.DeletedObjects) > 0 {
				log.Fatalf("Some objects failed to delete: %v", delRes.DeletedObjects)
			}

			totalDeleted += len(objects)

			// 更新marker
			marker = oss.Marker(lor.NextMarker)
			if !lor.IsTruncated {
				break
			}
		}
		return nil
	}
}
