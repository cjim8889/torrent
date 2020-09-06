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
	name string
	length int
	pieceLength int
	pieces [][]byte
	private int
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

	_ = processRawInfoName(&rawInfo, info)
	_ = processRawInfoLength(&rawInfo, info)
	_ = processRawInfoPrivate(&rawInfo, info)

	err := processRawInfoPieceLength(&rawInfo, info)
	if err != nil {
		return MetaInfoError{"illegal rawInfo: Invalid Piece Length"}
	}

	err = processRawInfoPieces(&rawInfo, info)
	if err != nil {
		return MetaInfoError{"illegal rawInfo: Invalid Pieces"}
	}

	return nil
}

func processRawInfoPrivate(rawInfo *map[string]bencode.BencodeCell, info *Info) error {
	rawPrivate, ok := (*rawInfo)["private"]
	if !ok {
		return MetaInfoError{"rawInfo does not contain private field"}
	}

	private, ok := rawPrivate.Value.(int)
	if !ok {
		return MetaInfoError{"rawInfor: private is not in valid form"}
	}

	info.private = private
	return nil
}

func processRawInfoPieces(rawInfo *map[string]bencode.BencodeCell, info *Info) error {
	rawPieces, ok := (*rawInfo)["pieces"]
	if !ok {
		return MetaInfoError{"rawInfo does not contain pieces field"}
	}

	pieces, ok := rawPieces.Value.(string)
	if !ok {
		return MetaInfoError{"rawInfo: pieces is not in valid form"}
	}

	piecesB := []byte(pieces)
	chunkCount := len(piecesB) / 20
	info.pieces = make([][]byte, chunkCount)

	for i := 0; i < chunkCount; i++ {
		tmp := make([]byte, 0, 20)
		info.pieces[i] = append(tmp, piecesB[i * 20 : i * 20 + 20]...)
	}

	return nil
}

func processRawInfoPieceLength(rawInfo *map[string]bencode.BencodeCell, info *Info) error {
	rawPieceLength, ok := (*rawInfo)["piece length"]
	if !ok {
		return MetaInfoError{"rawInfo does not contain piece field"}
	}

	pieceLength, ok := rawPieceLength.Value.(int)
	if !ok {
		return MetaInfoError{"rawInfo: piece length is not in valid form"}
	}

	info.pieceLength = pieceLength
	return nil
}

func processRawInfoLength(rawInfo *map[string]bencode.BencodeCell, info *Info) error {
	rawLength, ok := (*rawInfo)["length"]
	if !ok {
		return MetaInfoError{"rawInfo does not contain length field"}
	}

	length, ok := rawLength.Value.(int)
	if !ok {
		return MetaInfoError{"rawInfo: length is not in valid form"}
	}

	info.length = length
	return nil
}

func processRawInfoName(rawInfo *map[string]bencode.BencodeCell, info *Info) error {
	rawName, ok := (*rawInfo)["name"]
	if !ok {
		return MetaInfoError{"rawInfo does not contain name field"}
	}

	name, ok := rawName.Value.(string)
	if !ok {
		return MetaInfoError{"rawInfo: name is not in valid form"}
	}

	info.name = name
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

	_ = processMetaInfoInfo(&rawMeta, result)

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
