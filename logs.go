package main

func (spider *Spider) logErr(str string)  {
	TypeMonitorChan <- TypeErrNum
}

func (spider *Spider) logException(str string)  {
	TypeMonitorChan <- TypeExceptionNum
}

func (spider *Spider) logWarning(str string)  {
	TypeMonitorChan <- TypeWarningNum

}