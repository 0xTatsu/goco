package pkg

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

/*
	https://cloud.google.com/storage/docs/renaming-copying-moving-objects
	https://github.com/GoogleCloudPlatform/golang-samples/blob/master/storage/objects/main.go
*/

type GCS struct {
	bucketHandler *storage.BucketHandle
}

func NewGCS(
	storageClient *storage.Client,
	bucketName string,
) *GCS {
	return &GCS{
		bucketHandler: storageClient.Bucket(bucketName),
	}
}

func (gcs *GCS) Read(ctx context.Context, filePath string) ([]byte, error) {
	reader, err := gcs.bucketHandler.Object(filePath).NewReader(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = reader.Close(); err != nil {
			log.Printf("cannot close reader: %s", err)
		}
	}()

	fileByte, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return fileByte, nil
}

func (gcs *GCS) IsExistent(ctx context.Context, filePath string) bool {
	_, err := gcs.bucketHandler.Object(filePath).Attrs(ctx)
	return err != storage.ErrObjectNotExist
}

func (gcs *GCS) Copy(ctx context.Context, filePathSrc string, filePathDest string) error {
	src := gcs.bucketHandler.Object(filePathSrc)
	dst := gcs.bucketHandler.Object(filePathDest)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return err
	}

	return nil
}

func (gcs *GCS) Rename(ctx context.Context, oldFile string, newFile string) error {
	src := gcs.bucketHandler.Object(oldFile)
	dst := gcs.bucketHandler.Object(newFile)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return err
	}

	if err := src.Delete(ctx); err != nil {
		return err
	}

	return nil
}

func (gcs *GCS) Delete(ctx context.Context, filePath string) error {
	return gcs.bucketHandler.Object(filePath).Delete(ctx)
}

func (gcs *GCS) Update(ctx context.Context, filePath *os.File, newFilename string) error {
	wc := gcs.bucketHandler.Object(newFilename).NewWriter(ctx)
	if _, err := io.Copy(wc, filePath); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

func (gcs *GCS) makePublic(ctx context.Context, filePath string) error {
	acl := gcs.bucketHandler.Object(filePath).ACL()

	return acl.Set(ctx, storage.AllUsers, storage.RoleReader)
}

/*
	[START storage_list_files_with_prefix]
	Prefixes and delimiters can be used to emulate directory listings.
	Prefixes can be used filter objects starting with prefix.
	The delimiter argument can be used to restrict the results to only the
	objects in the given "directory". Without the delimiter, the entire  tree
	under the prefix is returned.

	For example, given these blobs:
	  /a/1.txt
	  /a/b/2.txt

	If you just specify prefix="a/", you'll get back:
	  /a/1.txt
	  /a/b/2.txt

	However, if you specify prefix="a/" and delim="/", you'll get back:
	  /a/1.txt
*/
func (gcs *GCS) ListByPrefixes(ctx context.Context, prefixes []string, delimiter string) ([]*storage.ObjectAttrs, error) {
	var objects []*storage.ObjectAttrs

	for _, prefix := range prefixes {
		objectIteration := gcs.bucketHandler.Objects(ctx, &storage.Query{
			Prefix:    prefix,
			Delimiter: delimiter,
		})

		for {
			objectAttrs, err := objectIteration.Next()
			if err == iterator.Done {
				break
			}

			if err != nil {
				return nil, errors.Wrapf(err, "cannot iterate bucket object of prefix: %s", prefix)
			}

			objects = append(objects, objectAttrs)
		}
	}

	return objects, nil
}
