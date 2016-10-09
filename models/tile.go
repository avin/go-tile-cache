package models

import (
	"os"
	"image/png"
	"github.com/harrydb/go/img/grayscale"
	"net/http"
	"io"
	"time"
	"path/filepath"
	"errors"
)

type TileManager struct {
	X string
	Y string
	Z string
	Server string
	GS bool
}

//create grayScale image version
func grayScale(filename string) string {

	//read file
	infile, err := os.Open(filename)
	if err != nil {
		panic(err.Error())
	}
	defer infile.Close()

	//decode PNG
	src, err := png.Decode(infile)
	if err != nil {
		panic(err.Error())
	}

	//convert
	gray := grayscale.Convert(src, grayscale.ToGrayLuminance)

	//get name for converted file
	outFileName := filename[:len(filename) - 4] + ".gs.png"

	if _, err := os.Stat(outFileName); err == nil {
		//if file exist - return it now
		return outFileName
	}

	//save converted file
	outfile, err := os.Create(outFileName)
	if err != nil {
		panic(err.Error())
	}
	defer outfile.Close()
	png.Encode(outfile, gray)

	return outFileName
}

//download file
func downloadFile(fileName string, url string) (error) {

	//get url file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if (resp.StatusCode != 200){
		return errors.New("error status code")
	}

	//create output file
	out, err := os.Create(fileName)
	if err != nil {
		panic(err.Error())
	}
	defer out.Close()

	//write content to output file
	io.Copy(out, resp.Body)

	return nil
}

// NewTaskManager returns an empty TaskManager.
func NewTileManager() *TileManager {
	return &TileManager{}
}

//get file from cache or download (and convert to GS)
func (tm *TileManager) Get(url string, ttl time.Duration) (string, error) {

	//make path
	path := filepath.Join("cache", tm.Server, tm.X, tm.Y)
	os.MkdirAll(path, os.ModePerm)
	fileName := filepath.Join(path, tm.Z + ".png")

	toDownload := false
	if file, err := os.Stat(fileName); err != nil {
		//if file not exist
		toDownload = true;
	} else {
		//if file older then ttl duration
		if (time.Now().Add(-ttl).After(file.ModTime())) {
			toDownload = true;
		}
	}

	//if we should download file
	if (toDownload){

		tries:=0 //tries to get file
		for {
			tries++
			err := downloadFile(fileName, url)
			//if file get success - exit try loop
			if (err == nil){
				break
			}

			//If cannot get file in 5 tries - return error
			if (err != nil && (tries>5)){
				return "", err
			}
		}
	}

	//if we gonna get grayScale tile - convert it
	if (tm.GS) {
		gsFileName := grayScale(fileName)
		return gsFileName, nil
	}

	return fileName, nil
}