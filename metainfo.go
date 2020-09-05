package main

import (
	"fmt"
	"github.com/cjim8889/bencode"
	"os"
)

type MetaInfo struct {
	announce string
	announceList []string
	info *Info
}

type Info struct {
	files []File
}

type File struct {
	length int
	path string
}

type MetaInfoError struct {
	error string
}

func (e MetaInfoError) Error() string {
	return fmt.Sprintf("MetaInfo Error: %v\n", e.error)
}

func processMetaInfoInfo(rawMeta *map[string]bencode.BencodeCell, result *MetaInfo) error {
	rawInfo, ok := (*rawMeta)["info"].Value.(map[string]bencode.BencodeCell)
	if !ok {
		return MetaInfoError{"illegal rawInfo"}
	}

	info := new(Info)
	result.info = info

	rawFiles, ok := rawInfo["files"].Value.([]bencode.BencodeCell)
	if !ok {
		return MetaInfoError{"illegal rawInfo: Files"}
	}

	files := make([]File, 0, len(rawFiles))
	for _, v := range rawFiles {
		rawFile := v.Value.(map[string]bencode.BencodeCell)

		length, ok := rawFile["length"]
		if !ok {
			return MetaInfoError{"illegal rawInfo: File: Length"}
		}

		path, ok := rawFile["path"]
		if !ok {
			return MetaInfoError{"illegal rawInfo: File: Path"}
		}

		files = append(files, File{length: length.Value.(int), path: path.Value.([]bencode.BencodeCell)[0].Value.(string)})
	}

	info.files = files

	return nil
}

func UnmarshalMetaInfoFrom(rawMeta map[string]bencode.BencodeCell) (*MetaInfo, error) {
	announce, ok := rawMeta["announce"]
	if !ok {
		return nil, MetaInfoError{"illegal rawMeta"}
	}

	result := new(MetaInfo)
	result.announce = announce.Value.(string)

	announceListRaw, ok := rawMeta["announce-list"]
	if ok {
		temp := announceListRaw.Value.([]bencode.BencodeCell)
		announceList := make([]string, 0, len(temp))
		for _, v := range temp {
			announceList = append(announceList, v.Value.([]bencode.BencodeCell)[0].Value.(string))
		}

		result.announceList = announceList
	}

	processMetaInfoInfo(&rawMeta, result)

	return result, nil
}



func main() {
	file, err := os.Open("/Users/wuhaochen/go/src/github.com/cjim8889/torrent/test.torrent")
	if err != nil {
		fmt.Println(err.Error())
	}


	bencodeReader := bencode.NewBencodeReader(file)
	metainfo, err := bencodeReader.DecodeStream()
	if err != nil {
		return
	}


	//fmt.Printf("%v\n", metainfo)
	r, err := UnmarshalMetaInfoFrom(metainfo.(map[string]bencode.BencodeCell))
	fmt.Println(r)
}
