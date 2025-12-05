package errors

import (
	"fmt"
	"strconv"
	"strings"
)

type ResultCode uint32

func (rc ResultCode) Failed() bool {
	return int32(rc) < 0
}

func (rc ResultCode) Error() string {
	str := rc.String()
	if strings.HasSuffix(str, ")") {
		return "0x" + strconv.FormatUint(uint64(rc), 16)
	}
	return str + " (0x" + strconv.FormatUint(uint64(rc), 16) + ")"
}

const (
	ResultSuccess                          ResultCode = 0x0
	ResultInvalidArg                       ResultCode = 0x80070057
	StatusOccluded                         ResultCode = 0x087A0001
	StatusClipped                          ResultCode = 0x087A0002
	StatusNoRedirection                    ResultCode = 0x087A0004
	StatusNoDesktopAccess                  ResultCode = 0x087A0005
	StatusGraphicsVidpnSourceInUse        ResultCode = 0x087A0006
	StatusModeChanged                      ResultCode = 0x087A0007
	StatusModeChangeInProgress             ResultCode = 0x087A0008
	ErrorInvalidCall                       ResultCode = 0x887A0001
	ErrorNotFound                          ResultCode = 0x887A0002
	ErrorMoreData                          ResultCode = 0x887A0003
	ErrorUnsupported                       ResultCode = 0x887A0004
	ErrorDeviceRemoved                     ResultCode = 0x887A0005
	ErrorDeviceHung                        ResultCode = 0x887A0006
	ErrorDeviceReset                       ResultCode = 0x887A0007
	ErrorWasStillDrawing                   ResultCode = 0x887A000A
	ErrorFrameStatisticsDisjoint           ResultCode = 0x887A000B
	ErrorGraphicsVidpnSourceInUse         ResultCode = 0x887A000C
	ErrorDriverInternalError               ResultCode = 0x887A0020
	ErrorNonexclusive                      ResultCode = 0x887A0021
	ErrorNotCurrentlyAvailable            ResultCode = 0x887A0022
	ErrorRemoteClientDisconnected          ResultCode = 0x887A0023
	ErrorRemoteOutOfMemory                 ResultCode = 0x887A0024
	ErrorAccessLost                        ResultCode = 0x887A0026
	ErrorWaitTimeout                       ResultCode = 0x887A0027
	ErrorSessionDisconnected               ResultCode = 0x887A0028
	ErrorRestrictToOutputStale             ResultCode = 0x887A0029
	ErrorCannotProtectContent              ResultCode = 0x887A002A
	ErrorAccessDenied                      ResultCode = 0x887A002B
	ErrorNameAlreadyExists                 ResultCode = 0x887A002C
	ErrorSdkComponentMissing               ResultCode = 0x887A002D
	ErrorNotCurrent                        ResultCode = 0x887A002E
	ErrorHwProtectionOutOfMemory           ResultCode = 0x887A0030
	ErrorDynamicCodePolicyViolation        ResultCode = 0x887A0031
	ErrorNonCompositedUi                   ResultCode = 0x887A0032
	StatusUnoccluded                       ResultCode = 0x087A0009
	StatusDdaWasStillDrawing               ResultCode = 0x087A000A
	ErrorModeChangeInProgress              ResultCode = 0x887A0025
	StatusPresentRequired                  ResultCode = 0x087A002F
	ErrorCacheCorrupt                      ResultCode = 0x887A0033
	ErrorCacheFull                         ResultCode = 0x887A0034
	ErrorCacheHashCollision                ResultCode = 0x887A0035
	ErrorAlreadyExists                     ResultCode = 0x887A0036
	DdiErrWasStillDrawing                  ResultCode = 0x887B0001
	DdiErrUnsupported                      ResultCode = 0x887B0002
	DdiErrNonexclusive                     ResultCode = 0x887B0003
)

func (rc ResultCode) String() string {
	switch rc {
	case ResultSuccess:
		return "ResultSuccess"
	case ResultInvalidArg:
		return "ResultInvalidArg"
	case StatusOccluded:
		return "StatusOccluded"
	case StatusClipped:
		return "StatusClipped"
	case StatusNoRedirection:
		return "StatusNoRedirection"
	case StatusNoDesktopAccess:
		return "StatusNoDesktopAccess"
	case StatusGraphicsVidpnSourceInUse:
		return "StatusGraphicsVidpnSourceInUse"
	case StatusModeChanged:
		return "StatusModeChanged"
	case StatusModeChangeInProgress:
		return "StatusModeChangeInProgress"
	case ErrorInvalidCall:
		return "ErrorInvalidCall"
	case ErrorNotFound:
		return "ErrorNotFound"
	case ErrorMoreData:
		return "ErrorMoreData"
	case ErrorUnsupported:
		return "ErrorUnsupported"
	case ErrorDeviceRemoved:
		return "ErrorDeviceRemoved"
	case ErrorDeviceHung:
		return "ErrorDeviceHung"
	case ErrorDeviceReset:
		return "ErrorDeviceReset"
	case ErrorWasStillDrawing:
		return "ErrorWasStillDrawing"
	case ErrorFrameStatisticsDisjoint:
		return "ErrorFrameStatisticsDisjoint"
	case ErrorGraphicsVidpnSourceInUse:
		return "ErrorGraphicsVidpnSourceInUse"
	case ErrorDriverInternalError:
		return "ErrorDriverInternalError"
	case ErrorNonexclusive:
		return "ErrorNonexclusive"
	case ErrorNotCurrentlyAvailable:
		return "ErrorNotCurrentlyAvailable"
	case ErrorRemoteClientDisconnected:
		return "ErrorRemoteClientDisconnected"
	case ErrorRemoteOutOfMemory:
		return "ErrorRemoteOutOfMemory"
	case ErrorAccessLost:
		return "ErrorAccessLost"
	case ErrorWaitTimeout:
		return "ErrorWaitTimeout"
	case ErrorSessionDisconnected:
		return "ErrorSessionDisconnected"
	case ErrorRestrictToOutputStale:
		return "ErrorRestrictToOutputStale"
	case ErrorCannotProtectContent:
		return "ErrorCannotProtectContent"
	case ErrorAccessDenied:
		return "ErrorAccessDenied"
	case ErrorNameAlreadyExists:
		return "ErrorNameAlreadyExists"
	case ErrorSdkComponentMissing:
		return "ErrorSdkComponentMissing"
	case ErrorNotCurrent:
		return "ErrorNotCurrent"
	case ErrorHwProtectionOutOfMemory:
		return "ErrorHwProtectionOutOfMemory"
	case ErrorDynamicCodePolicyViolation:
		return "ErrorDynamicCodePolicyViolation"
	case ErrorNonCompositedUi:
		return "ErrorNonCompositedUi"
	case StatusUnoccluded:
		return "StatusUnoccluded"
	case StatusDdaWasStillDrawing:
		return "StatusDdaWasStillDrawing"
	case ErrorModeChangeInProgress:
		return "ErrorModeChangeInProgress"
	case StatusPresentRequired:
		return "StatusPresentRequired"
	case ErrorCacheCorrupt:
		return "ErrorCacheCorrupt"
	case ErrorCacheFull:
		return "ErrorCacheFull"
	case ErrorCacheHashCollision:
		return "ErrorCacheHashCollision"
	case ErrorAlreadyExists:
		return "ErrorAlreadyExists"
	case DdiErrWasStillDrawing:
		return "DdiErrWasStillDrawing"
	case DdiErrUnsupported:
		return "DdiErrUnsupported"
	case DdiErrNonexclusive:
		return "DdiErrNonexclusive"
	default:
		return fmt.Sprintf("UnknownResultCode(0x%x)", uint32(rc))
	}
}

