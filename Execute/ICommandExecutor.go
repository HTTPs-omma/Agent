package Execute

type ICommandExecutor interface {
	execute(cmd string) (string, error) // 실행
}
