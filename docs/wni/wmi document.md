refer : https://pkg.go.dev/github.com/stackexchange/wmi

## Overview
- wmi 패키지는 윈도우의 WMI 과 상호작용 하는 WQL 을 제공합니다.
- 윈도우 환경에서만 동작한다는 것을 명심해야 합니다.
- example : 
```go
type Win32_Process struct {
	Name string
}

func main() {
	var dst []Win32_Process
	q := wmi.CreateQuery(&dst, "")
	err := wmi.Query(q, &dst)
	if err != nil {
		log.Fatal(err)
	}
	for i, v := range dst {
		println(i, v.Name)
	}
}
```

#### func CallMethod(connectServerArgs []interface{}, className, methodName string, params []interface{}) (int32, error)
- CallMethod는 주어진 params를 사용하여 className이라는 클래스의 인스턴스에서 methodName이라는 이름의 메서드를 호출합니다.

#### func CreateQuery(src interface{}, where string, class ...string) string
- CreateQuery 함수는 src의 모든 열을 쿼리하는 WQL 쿼리 문자열을 반환합니다. 
- where는 선택적 문자열로, WHERE 절과 함께 사용될 때 쿼리에 추가됩니다. 
- 이 경우 "WHERE" 문자열은 맨 앞에 위치해야 합니다. wmi 클래스는 타입의 이름으로 얻어집니다. 
- 익명 구조체에 유용한 선택적 클래스 매개변수를 통해 클래스를 전달할 수 있습니다.

#### func Query(query string, dst interface{}, connectServerArgs ...interface{}) error
- Query 함수는 WQL 쿼리를 실행하고 결과 값을 dst에 추가합니다

## WQL(WMI Query Language)
- WQL은 WMI(Windows Management Instrumentation) 쿼리 언어로, WMI에서 정보를 가져오는 데 사용되는 언어입니다.
- refer : https://learn.microsoft.com/ko-kr/powershell/module/microsoft.powershell.core/about/about_wql?view=powershell-5.1
- WQL 쿼리는 표준 Get-WmiObject 명령보다 다소 빠르며 명령이 수백 개의 시스템에서 실행될 때 향상된 성능이 분명
- WMI, 특히 WQL을 사용하는 경우 Windows PowerShell도 사용 중이라는 사실을 잊지 마세요. WQL 쿼리가 예상대로 작동하지 않는 경우 WQL 쿼리를 디버그하는 것보다 표준 Windows PowerShell 명령을 사용하는 것이 더 쉽습니다.

- SELECT <property> FROM <WMI-class>
WMI Class list : https://learn.microsoft.com/ko-kr/windows/win32/wmisdk/wmi-classes
- 


