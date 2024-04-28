package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
)

func compressJSON(data interface{}) ([]byte, error) {
	// 将数据编码为 JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// 创建一个 buffer 用于保存压缩后的数据
	var compressedData bytes.Buffer

	// 创建一个 gzip writer
	gzipWriter := gzip.NewWriter(&compressedData)

	// 将 JSON 数据写入 gzip writer，并进行压缩
	_, err = gzipWriter.Write(jsonData)
	if err != nil {
		return nil, err
	}

	// 关闭 gzip writer，确保所有数据都被写入
	err = gzipWriter.Close()
	if err != nil {
		return nil, err
	}

	// 返回压缩后的数据
	return compressedData.Bytes(), nil
}

func decompressJSON(compressedData []byte) (string, error) {
	// 创建一个 buffer 用于保存解压后的数据
	var decompressedData bytes.Buffer

	// 创建一个 bytes.Reader 来读取压缩的数据
	compressedDataReader := bytes.NewReader(compressedData)

	// 创建一个 gzip reader，并从压缩的数据中读取并解压缩
	gzipReader, err := gzip.NewReader(compressedDataReader)
	if err != nil {
		return "", err
	}
	defer gzipReader.Close()

	// 将解压后的数据写入 buffer
	_, err = decompressedData.ReadFrom(gzipReader)
	if err != nil {
		return "", err
	}

	// 解压缩后的数据是 JSON 格式，将其解码为 interface{} 类型
	var jsonData string
	err = json.Unmarshal(decompressedData.Bytes(), &jsonData)
	if err != nil {
		return "", err
	}

	return jsonData, nil
}

func main() {
	// 原始 JSON 数据
	jsonData := `{"company_code":"91370900676823234F","collecttime":"20240425085608","indicator_code":"5gufengD504.PV","indicator_value":"576.900024","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F"collecttime":"20240425085608","indicator_code":"5gufengD506.PV","indicator_value":"51723148.000000","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085623","indicator_code":"5gufengD504.PV","indicator_value":"579.299988","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085623","indicator_code":"5gufengD506.PV","indicator_value":"51723152.000000","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085638","indicator_code":"5gufengD504.PV","indicator_value":"579.299988","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085638","indicator_code":"5gufengD506.PV","indicator_value":"51723152.000000","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085653","indicator_code":"5gufengD504.PV","indicator_value":"586.099976","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085653","indicator_code":"5gufengD506.PV","indicator_value":"51723152.000000","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085708","indicator_code":"5gufengD504.PV","indicator_value":"586.099976","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085708","indicator_code":"5gufengD506.PV","indicator_value":"51723152.000000","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085723","indicator_code":"5gufengD504.PV","indicator_value":"590.500000","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085723","indicator_code":"5gufengD506.PV","indicator_value":"51723156.000000","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085738","indicator_code":"5gufengD504.PV","indicator_value":"607.400024","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085738","indicator_code":"5gufengD506.PV","indicator_value":"51723160.000000","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085753","indicator_code":"5gufengD504.PV","indicator_value":"633.799988","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085753","indicator_code":"5gufengD506.PV","indicator_value":"51723164.000000","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085808","indicator_code":"5gufengD504.PV","indicator_value":"648.799988","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085808","indicator_code":"5gufengD506.PV","indicator_value":"51723168.000000","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085823","indicator_code":"5gufengD504.PV","indicator_value":"648.799988","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085823","indicator_code":"5gufengD506.PV","indicator_value":"51723168.000000","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085838","indicator_code":"5gufengD504.PV","indicator_value":"690.700012","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085838","indicator_code":"5gufengD506.PV","indicator_value":"51723172.000000","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085853","indicator_code":"5gufengD504.PV","indicator_value":"690.700012","device_code":"3fef34a75c3c66f0737150b72743e630"},{"company_code":"91370900676823234F","collecttime":"20240425085853","indicator_code":"5gufengD506.PV","indicator_value":"51723172.000000","device_code":"3fef34a75c3c66f0737150b72743e630"}]}`

	// 压缩 JSON 数据
	compressedJSON, err := compressJSON(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	// 打印压缩后的数据长度和内容
	fmt.Printf("Compressed JSON size: %d bytes\n", len(compressedJSON))
	fmt.Printf("Compressed JSON: %s\n", string(compressedJSON))

	data, err := decompressJSON(compressedJSON)

	fmt.Println("from  JSON:", data, err)
}
