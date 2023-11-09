package s3

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sarulabs/di/v2"
)

var Definitions = []di.Def{
	{
		Build: func(ctn di.Container) (interface{}, error) {
			client := ctn.Get("focus.awsS3.client").(*s3.Client)
			bucket := ctn.Get("focus.awsS3.bucketName").(string)
			return NewMediaStorage(client, bucket)
		},
		Name: "focus.models.fileStorage",
	},
}
