package utils

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

func ReplaceImage(filePath, newImage string) error {
	// ファイルを読み込む
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// 正規表現パターンを定義
	// `image:`の後に続く任意の文字列をキャプチャ
	re := regexp.MustCompile(`(?m)^(\s*image:\s*)(.+)$`)

	// 新しいイメージに置き換える
	result := re.ReplaceAllString(string(data), fmt.Sprintf("${1}%s", newImage))

	// ファイルに書き込む
	err = ioutil.WriteFile(filePath, []byte(result), 0644)
	if err != nil {
		return err
	}

	return nil
}
