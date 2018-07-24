package redis

import (
	"bytes"

	"github.com/tatsushid/go-critbit"
)

// reqType is a kind for cmd which is Ctl/Read/Write/NotSupport
type reqType = int

const (
	reqTypeCtl        reqType = iota
	reqTypeRead       reqType = iota
	reqTypeWrite      reqType = iota
	reqTypeNotSupport reqType = iota
)

var reqMap *critbit.Tree

var (
	reqReadBytes = []byte("" +
		"4\r\nDUMP" +
		"6\r\nEXISTS" +
		"4\r\nPTTL" +
		"3\r\nTTL" +
		"4\r\nTYPE" +
		"8\r\nBITCOUNT" +
		"6\r\nBITPOS" +
		"3\r\nGET" +
		"6\r\nGETBIT" +
		"8\r\nGETRANGE" +
		"4\r\nMGET" +
		"6\r\nSTRLEN" +
		"7\r\nHEXISTS" +
		"4\r\nHGET" +
		"7\r\nHGETALL" +
		"5\r\nHKEYS" +
		"4\r\nHLEN" +
		"5\r\nHMGET" +
		"7\r\nHSTRLEN" +
		"5\r\nHVALS" +
		"5\r\nHSCAN" +
		"5\r\nSCARD" +
		"5\r\nSDIFF" +
		"6\r\nSINTER" +
		"9\r\nSISMEMBER" +
		"8\r\nSMEMBERS" +
		"11\r\nSRANDMEMBER" +
		"6\r\nSUNION" +
		"5\r\nSSCAN" +
		"5\r\nZCARD" +
		"6\r\nZCOUNT" +
		"9\r\nZLEXCOUNT" +
		"6\r\nZRANGE" +
		"11\r\nZRANGEBYLEX" +
		"13\r\nZRANGEBYSCORE" +
		"5\r\nZRANK" +
		"9\r\nZREVRANGE" +
		"14\r\nZREVRANGEBYLEX" +
		"16\r\nZREVRANGEBYSCORE" +
		"8\r\nZREVRANK" +
		"6\r\nZSCORE" +
		"5\r\nZSCAN" +
		"6\r\nLINDEX" +
		"4\r\nLLEN" +
		"6\r\nLRANGE" +
		"7\r\nPFCOUNT")

	reqWriteBytes = []byte("" +
		"3\r\nDEL" +
		"6\r\nEXPIRE" +
		"8\r\nEXPIREAT" +
		"7\r\nPERSIST" +
		"7\r\nPEXPIRE" +
		"9\r\nPEXPIREAT" +
		"7\r\nRESTORE" +
		"4\r\nSORT" +
		"6\r\nAPPEND" +
		"4\r\nDECR" +
		"6\r\nDECRBY" +
		"6\r\nGETSET" +
		"4\r\nINCR" +
		"6\r\nINCRBY" +
		"11\r\nINCRBYFLOAT" +
		"4\r\nMSET" +
		"6\r\nPSETEX" +
		"3\r\nSET" +
		"6\r\nSETBIT" +
		"5\r\nSETEX" +
		"5\r\nSETNX" +
		"8\r\nSETRANGE" +
		"4\r\nHDEL" +
		"7\r\nHINCRBY" +
		"12\r\nHINCRBYFLOAT" +
		"5\r\nHMSET" +
		"4\r\nHSET" +
		"6\r\nHSETNX" +
		"7\r\nLINSERT" +
		"4\r\nLPOP" +
		"5\r\nLPUSH" +
		"6\r\nLPUSHX" +
		"4\r\nLREM" +
		"4\r\nLSET" +
		"5\r\nLTRIM" +
		"4\r\nRPOP" +
		"9\r\nRPOPLPUSH" +
		"5\r\nRPUSH" +
		"6\r\nRPUSHX" +
		"4\r\nSADD" +
		"10\r\nSDIFFSTORE" +
		"11\r\nSINTERSTORE" +
		"5\r\nSMOVE" +
		"4\r\nSPOP" +
		"4\r\nSREM" +
		"11\r\nSUNIONSTORE" +
		"4\r\nZADD" +
		"7\r\nZINCRBY" +
		"11\r\nZINTERSTORE" +
		"4\r\nZREM" +
		"14\r\nZREMRANGEBYLEX" +
		"15\r\nZREMRANGEBYRANK" +
		"16\r\nZREMRANGEBYSCORE" +
		"11\r\nZUNIONSTORE" +
		"5\r\nPFADD" +
		"7\r\nPFMERGE")

	reqNotSupportBytes = []byte("" +
		"6\r\nMSETNX" +
		"5\r\nBLPOP" +
		"5\r\nBRPOP" +
		"10\r\nBRPOPLPUSH" +
		"4\r\nKEYS" +
		"7\r\nMIGRATE" +
		"4\r\nMOVE" +
		"6\r\nOBJECT" +
		"9\r\nRANDOMKEY" +
		"6\r\nRENAME" +
		"8\r\nRENAMENX" +
		"4\r\nSCAN" +
		"4\r\nWAIT" +
		"5\r\nBITOP" +
		"4\r\nEVAL" +
		"7\r\nEVALSHA" +
		"4\r\nAUTH" +
		"4\r\nECHO" +
		"4\r\nINFO" +
		"5\r\nPROXY" +
		"7\r\nSLOWLOG" +
		"4\r\nQUIT" +
		"6\r\nSELECT" +
		"4\r\nTIME" +
		"6\r\nCONFIG" +
		"8\r\nCOMMANDS")

	reqCtlBytes = []byte("4\r\nPING")
)

func getReqType(cmd []byte) reqType {
	if bytes.Contains(reqNotSupportBytes, cmd) {
		return reqTypeNotSupport
	}

	if bytes.Contains(reqReadBytes, cmd) {
		return reqTypeRead
	}

	if bytes.Contains(reqWriteBytes, cmd) {
		return reqTypeWrite
	}

	if bytes.Contains(reqCtlBytes, cmd) {
		return reqTypeCtl
	}

	return reqTypeNotSupport
}
