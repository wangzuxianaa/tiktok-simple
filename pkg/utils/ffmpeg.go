package utils

import (
	"github.com/disintegration/imaging"
	"github.com/u2takey/ffmpeg-go/examples"
)

//
// ReadFrameAsJpeg
// @Description: 读取视频一帧作为图片
// @param inFileName 视频流的地址
// @param frameNum 哪一帧图片
// @param outFileName 图片保存地址
// @return error
//
func ReadFrameAsJpeg(inFileName string, frameNum int, outFileName string) error {
	reader := examples.ExampleReadFrameAsJpeg(inFileName, frameNum)
	img, err := imaging.Decode(reader)
	if err != nil {
		return err
	}
	// 图片保存到指定地址
	err = imaging.Save(img, outFileName)
	if err != nil {
		return err
	}
	return nil
}
