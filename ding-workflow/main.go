// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dingtalkworkflow_1_0 "github.com/alibabacloud-go/dingtalk/workflow_1_0"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	"os"
)

/**
 * 使用 Token 初始化账号Client
 * @return Client
 * @throws Exception
 */
func CreateClient() (_result *dingtalkworkflow_1_0.Client, _err error) {
	config := &openapi.Config{}
	config.Protocol = tea.String("https")
	config.RegionId = tea.String("central")
	_result = &dingtalkworkflow_1_0.Client{}
	_result, _err = dingtalkworkflow_1_0.NewClient(config)
	return _result, _err
}

func _main(args []*string) (_err error) {
	name := "数字输入框"
	value := "1"
	client, _err := CreateClient()
	if _err != nil {
		return _err
	}

	startProcessInstanceHeaders := &dingtalkworkflow_1_0.StartProcessInstanceHeaders{}
	startProcessInstanceHeaders.XAcsDingtalkAccessToken = tea.String("dbccd9f791c9393a9661cfb3dcac0341")
	startProcessInstanceRequest := &dingtalkworkflow_1_0.StartProcessInstanceRequest{
		OriginatorUserId: tea.String("22630306521225782"),
		ProcessCode:      tea.String("PROC-1179F945-9C6B-49E5-9A72-F7C9D8008DB2"),
		CcPosition:       tea.String("START"),
		FormComponentValues: []*dingtalkworkflow_1_0.StartProcessInstanceRequestFormComponentValues{
			{
				Name:  &name,
				Value: &value,
			},
		},
	}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		_, _err = client.StartProcessInstanceWithOptions(startProcessInstanceRequest, startProcessInstanceHeaders, &util.RuntimeOptions{})
		if _err != nil {
			return _err
		}

		return nil
	}()

	if tryErr != nil {
		var err = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			err = _t
			fmt.Println(_t)
		} else {
			err.Message = tea.String(tryErr.Error())
		}
		if !tea.BoolValue(util.Empty(err.Code)) && !tea.BoolValue(util.Empty(err.Message)) {

			return _err
		}

	}
	return _err
}

var data = `
SDKError:
   StatusCode: 400
   Code: InvalidAuthentication
   Message: code: 400, 不合法的access_token request id: 03F2E3EB-1F9C-77EA-B8D2-7DBD6DA4671E
   Data: {"code":"InvalidAuthentication","message":"不合法的access_token","requestid":"03F2E3EB-1F9C-77EA-B8D2-7DBD6DA4671E","statusCode":400}

`

func main() {
	err := _main(tea.StringSlice(os.Args[1:]))
	if err != nil {
		log.Fatal(err)
	}
}
