package main

import (
	"fmt"
	"strings"
	"unicode"
)

func main() {
	fmt.Println(part1(input1))
	fmt.Println(part1(input2))
	fmt.Println(part2(input1))
	fmt.Println(part2(input2))
}

func part1(input string) int {
	var total int
	for _, row := range strings.Split(input, "\n") {
		mid := len(row) / 2
		items := overlap(row[:mid], row[mid:])
		for _, item := range items {
			total += priority(item)
		}

	}
	return total
}

func part2(input string) int {
	var total int
	rows := strings.Split(input, "\n")
	for i := 0; i < len(rows); i += 3 {
		l, m, r := rows[i], rows[i+1], rows[i+2]
		items := string(overlap(string(overlap(l, m)), r))
		for _, item := range items {
			total += priority(item)
		}

	}
	return total
}

func priority(r rune) int {
	if r == unicode.ToLower(r) {
		return int(r-'a') + 1
	}
	return int(r-'A') + 1 + 26
}

func overlap(a, b string, rest ...string) []rune {
	if len(a) < len(b) {
		return overlap(b, a)
	}
	cache := make(map[rune]bool)
	for _, c := range a {
		cache[c] = true
	}
	var result []rune
	found := make(map[rune]bool)
	for _, c := range b {
		if cache[c] && !found[c] {
			found[c] = true
			result = append(result, c)
		}
	}
	return result
}

var input1 = `vJrwpWtwJgWrhcsFMMfFFhFp
jqHRNqRjqzjGDLGLrsFMfFZSrLrFZsSL
PmmdzqPrVvPwwTWBwg
wMqvLMZHhHMvwLHjbvcjnnSBnvTQFn
ttgJtRGJQctTZtZT
CrZsJsPPZsGzwwsLwLmpwMDw`

var input2 = `BdbzzddChsWrRFbzBrszbhWMLNJHLLLLHZtSLglFNZHLJH
nnfMwqpQTMffHlNNLllHnZSS
cGpcMwfppfqcjcTCBBzWDsDbDrjzWz
LhfjhcdjcGdhFfdGfdjdvwCCZMvvLvWwMLCLSwZC
rDnsbmptPmlbQMCrQWQQBZQW
gltgVPngDPbptPsbPzVgmDldfTdfczThjJJjfMcJdFHjjH
dGlgDflTLLLrRLTLVdQLcQMnbvHbbFzNNvMbnHHn
sZjWJJCSjWqfCqSjSmJSbFvCzBMBBzHncHNvMBHN
twqqwpZwfrlwRwDGDR
zCGGFTQMQrsNRNGZdR
cLLQgPDpgcgmvPRHrwBdvrNZ
glWpmJWQDcJpQnpjSmbhFtMnqFfCVTCFCFFM
zNZWFNZBFrGTdBcZZBdJTrGrmgvppgDHwHmgVHCHCvCPDjzC
qtqqPnLSfLwvjvvgvvqH
MtbLLLQbRfPRfnbnnLMtnsbdBNsrGWNWcNcTBWZrrcrWcJ
sZwstbbDVlHtbhcrhhZLrRpNQN
jqGjjFjWnzWGgqWjJJNphnLMRhLhcrhcrSLN
qCJzFJdvmHvbtpbb
ZSRsvvQvZpsRQGJghClPCwGPChCP
FVdMLDdtDRdDcBtmcVFntwgJJTlnNPwJnlTlwlTPgN
VqttMWFmDbjbzrSWRQ
TsDSBcwshdwSCrgRWZBvgGRG
LPVJLqqJbbzpFqwpbvgGRmZPrmZgCvCCfr
tzJlJzQllFLqtwHhjNSdtthjDhTN
fsDLDDnwvnSdqLSsFSDfsLpbgVttPMpPNjMWVMNfpjgW
jhHmBmlrTBBHRPVtMZbppNPPZB
JmTTrTmjRTJqSQQSFqQw
HPzZFgPFMCHJCcZMcDQGwpLqPLqppwhGvv
BrWRqbqRsTSTqNddrVrthLQVwnpwphnvnDnGvGDn
tSbBfsRbTfqjJqjCqm
pCqrqzmqZzrmCCvCJwJPBRwJBWBmwWBJ
VqbqbjFLFfSHnfctBwDdDFTwtRRTDF
LVcnbjHLSqHnhbSGGppCMMZMphpNrQ
PhTcTPsSPCMvvhhMRPttbNWfNsWFNfWWtpNw
rdQrDbJBVVjrBVdLjHHHWNWwfHtzzNtFpZZptppg
JJVGGdddjDjDJmdjGqqRSbPMTvcRlqnnMlvT
SqGfTrBlSrrrfGGQvCnqZhZmPPhvJh
LdVNwgsjdjHmjWwDsDpwsHWtnCnQQQQNnvCnbJZQbNbtZZ
DjssDHLLVppDssdLspswFLVjzfFMTfSBcTRcrFSBzBTzmGzr
JSJJqlldpJlcdVWMlgMJrcCjrhzHCwTjHrQzwTzZ
bBvsGBGNFDNRFNFBRDPsDDHCChhhrZRQQzhjhCjTzWCr
FvbLFGDDtWDFBnPGFDWGqJMgJpSdVStllttlppVS
clpBdBQBsqGpQbVdqTTWRTSFgLLggffg
NzvwmHvZtZZgbDSCLDmfmm
jPJjthHzNwPvvjwjNzPzztZBlrccpbVGBQhlGBVlrpnpBs
NmFFGlGmzCrNWHvFmFWCFvQPTdDDlbbdgJPtgbbdPDcc
RwfBqLwZqffqnPbdDgVQDdtRbT
hnptqLMMLwqjMMjzSzHmhvNNrvvhhz
hqVrDdPdVDqjsDrjjqsfrrWlctvHJNLfvNcLHNNZRNZHvfZL
CSnWQSGBBSnmBnTmSbQbNZZMJNJMccGNHLZNNvNv
gFnmBmBnTCpwmBbgWVqjPsrsjplVrjlqPr
vgVgJJCphzFjzbwljwww
PWmfDgrPrPWlSWqSWlSZ
rrTHQTGTPHDDHgPrcrcPmDPhdtGvnvNnttJCdJhtGVnpdM
DDDhNgWNLWgDqDgtgtwSngjVSQdf
cvFrcGmBrrrCdSfSJQ
mBHzFFvFzmsBspsFsZqhqbWlWdPDlLHbqZ
TLNpGpRzwGQLLQRTwdvWdWbdbgdBlblb
FJDVzZzHfZHVzcHgnvHnvngvvc
MzDVSmZPPrqhGmCqQqLR
mqHWcBHVcgsbhhnTrrTg
fWftGtfJpwJMSdFDLFSGDGFnnNrbrhTNsZnrhTswZZnnsb
dSftppdSFStLDpQfLjWHQvWBmmqVjCqcVv
frfNzgvzzzJwJqpRcP
VdVSnGnHqhDDTdqhdLWmjMTmjMPjPmTsmsmjmR
BdGLSWtBSWDtVqdSDVGqtSSHFZFZtQbZlZfNgFffrgNrvgrl
jBVSjsJcLcqqjtgcmRFRNFFzFm
CnHnWQQGGWnnCnfvCTmZRZgNtfZtNDZPtNDtzl
mWWGGdmQQWwrHHMBBBhrBVqshBLBSL
cLtFcllvrslGLcLHVzDZQzzpznWzQtQZ
TmShfSSPJRRBBfSgmdfBJhDNVpGzVZpzQgZgbWzpnDGp
BhJPfSqdhqJRdBdPqTJBChvrFlHvvCGCccHsMscvrHvv
zfddZTpZLzLDfLtDCttdTfZPnlcPcnhjPDnlcjMchDPjnP
MQsFvFHJsQvmNvvswljgbPbwlwjnmcch
VFVqvqFQHqVJRVCBLpqTGMzCqtMZ
PZdVgNdQQcdcZQtGhnHtBlvlvWTnbBHWWH
CDJmzqFmMmLfqmzFwqfzfwMRvpRbTWBBWDnlnbnWnTNTnB
zFfffjCqwLJCfrCFjCLCzCMsjSZjQNPhQGVSQQZhsGjGgd
nbHntnqPQwTHwQVC
BzfSZSpcBZpzpPhSBjRTCBWTTRWTTTWR
pvZpzPzNfhddJGmmNnJb
GwTgWlvbgTwMrbwTrlWvwtHBNtNvBhBtdZcShHDtNS
PnndmnCmnJFnsmRdmFnnZDhRRScRQQHDtDNtcDQS
CLffLsmqqpJljVdlpMpWlr
dPCzBLTSLqmqdSCsmrTDVQjZfjfVVZnZhhhLGQ
wwFpgvPJNwgPPwvZGnQhbZQQFfQZbf
WvRwwJpHHgpNMNqdTdCdqrmHBPtm
DbWwjSGFSFfwGfCwDSSPPjLhgrrLWRJRgggJphLzpJLq
lQTnMHdcQBvlHMMZBcHHTrbzpdzgqrpLdqzzVRVRVV
vlvlNlBQvBZTQBnHnQTTBBPwPNCCsfSFmbSFfmGPwFmj
hGGQtbVjhRqlmqqrmDlpmw
gPMZsMgdssCPPsvrgZcTZTPSnnLLnBWDwLmwWwBnmWSnWNWB
TsvgMPdcgCfMdcsZJRQhVfVRVQFFQGtr
NfpFTTpFNbpZMRFrJMMMCv
dWJPngDWBtPVBdPVGHZzCGHZrrvZRSMSzv
WDDDVDPlWnLBVgnsgJmQJNqThNmbjlqbbh
vnznqvfrzzVzrvvnfVqztBtGbMCdGmCcdccJccCLCcpSSgcL
RQsDsljDlhssWshHhlhsHTlLbgJLpmgbbcMSpSSbcSJgjd
WswhhHWlRTsQDwWHTRhsvNVvwqzfVmNBtZNmnzzq
cjcPfLlQtPsZQlfHZJqVSFdVwmSRRqSSddwDvF
MgNbBgzgnwdPRFmSPn
CMGbNNMhCMzzPzBpTNPGBclLcLfsptHQfQlZssLcLf
slsdfpSlllpTVJJGgGDgHMdV
wrBQrbQrhQcpbQrhLwRBrjVVgGGPgZMtZMVhgMPMGHPg
QQwRnrwRbbRcmQmrjRnjpvNsNTNSlTmsTqCSvsNSWz
GWNwppdHdpmzgPFPCRmlCBPB
bSrJhJSsMhrJMDPRCPBlwVCtVSLV
QbsbwqQZvrJhhQrhZZrhchTfnTWGNvWpNnjHGvGjNdGT
NMZGmnMBWmwmNnGwHrHvHzfrvrVVVj
pSbDRLgbpJDPpRZRQjjFqVhhDFjDVqfrzv
sQQscLTZcsTpRsBnmlBcdCMBMtlc
hwWslbGWbRvLZvcscZ
gQnmmrNTmSnTfgwDwVwpJvJPJzLqLJTLLvRcZz
mrmnSrDNNQSmmwdggrrDbHMGhtdtGhGlhtClMGhM
qQdlGcvDQDQvDdmtPmmmlStbjSrm
CpNhzWTCTRznBMvwtTjMrHtSvj
nWzsfZCsBhNpZLgGdQDddvqd
ggjTjJWDVVVRTwQcZWvchshWhs
LmFfLfSmBtCttNLfCFBPFNBvvvhrcQdvsrsdSwdqcwrwSw
CHttltmlPLMHRTgJQDgRDTMb
RPJgCdhgPPSzvWDcCfGHDcvf
HbrrwBspTwWDDnqbZjjD
rQrFsrsstQMQHJdHQm
GVwQVGBZBNQwsjdNcMMlgJNPgj
SWFfSzTCSWFCSpgnJLSjpMdc
FhTbvrzrMrDCVHsVsHGBtHwv
FsqjjVzFVWFqRRWBssdpNSBHwJpHHJJdddSN
lQgmhvbTcgTgfhTQhSCFCGJHSlwNtJdGHG
vgZPbbfMhbTmchhjRFRnsDRPWVqzWz
cvwfjjcJjqhctvSpCgCFVhrFCrpC
mRsQmsMlNNzznWQlRnsMRQlSCpLNbpwSSSrwFLDFLLLgFp
zznlnGlmRmMlPZnGQzMszMRfccTjcBJjtJtJjqTZtfwBTt
PtCwCCVqbjNNqqvGssPmHGsHMfPH
dcddcWFDJJJcLczhWQdcDScZvHZgpGfsvMfSmsggSvvnvM
TTQmTDhdQQzdRwjBVrRwCwBbBr
BnBsFHCrcnBrMBPSmMSCmrcFqnZLddLhdhdGLvqLqgqnLJJp
WfjTlNVDTtjzNWTlVMWNlTwTJphdgpvJLGLdgghZvGLpwddq
zVjbVblVDlNTlbzTzDjllDTzHmHCsmcMHcSBrCBbrHBcBPSC
ddlcGQlCjQNGQmPLslZTlmTtfT
MDzMwSwqpzpDRpWRJwgZhttrZmTPfZmrmtrMZZ
JSqRBpJzwgDDpqDqvpBRdCjbCjQCFCbHvdGPjdHQ
bwzPwGLZMsbJMPPLGLMQzbhhQRvWWtVqVhgCVtWWQDqt
HddrHFnFNpVnVLhnvLRV
NLdrFjHTBrrdjFSpFmNBmSfZMwZMJJfSffbwzwPbczbJ
QTWSzTTLwTfwflSNJRdvGlRGcNgJBl
FnmmqrqbBBgRbHGc
MFZqrCVCqmZCprspFZmnnMsDfzWzDwSfjSwPQPTLhffwBLwj
npfgFRTZRRnDZLdgRfRrrjcWzGpWGGGQrjjWpP
vblVbvSShhVzHWjzzlrPWG
bVwqvCBtShqBhCCtqhNqCJRTTBDFJJLnJgDRWFBnWg
nHDNQvgvnNZHDnsGcjfNTrTfVrfL
SRWFSBRLtSFqjTrVVcsVjTSG
BdBbRttWBdbdWdbppmZlLmgwHvgQvgLZ
PQRZlpDDptQSclBMGVBdhVFGBMpf
nnrsTCWjLJsnsSMShhGJfwffwV
vjzqsjqgSHbbtvHZ
DgFmbdSDZbPgLbDDmFwZwgLSfccGcGvnvvngsGGnsGMNWsWs
HqztHHhHVhHjhRRhJtCVBCfNprMWpcMMJfsvvMsGscpG
tjVtBVTBtHHqCRqtzQwTPSNdSwPTTbPFDdml
sbmBmHZbRRRmwBmsSjHzRjmSCNFLNLLQNQhFgtQLzNztlrff
MpqPPDPVnGqrJpcqqJpMVlgtlLFLChCgChCCQgCtCD
VdcVVJvdrVWHbwHWBSSBRb
tPDVBzzNSNdDSQpMQMTQJMMQMN
LqSWSmbsmfQTTGZMWGCW
cLjLcbrjqmvfqLbfmqLwDBBzSHPzlwzcdBlnnP
SbnHrGHSrrhHJBrrScDfcPDMfpPGcGcpDL
QTpmpmmQWlZsTvVQDMgggFLgFcPf
zsCNlltCslzlTNsNlShwdJCpSqdHrBhpwr
JZmFrmLGjFZdDGrrVTvzmPVvRRVzwzzl
WBnfMDBqMsgDBqpBvzwVlCwRTRRPpvlw
WfghfttggfSnnqbDftfqSBBDQdHhZFFJrrHFjGJdGjFrjQLj
rNLRjrlVlrFRJzlsjrVlRFGCmnMtftgCNDDgDmCfqNgPNM
SpdBpdHbpHWhZqnCdtCGggqCPn
QHHvvWwWhwVVQRscVzGl
QffLtMQGMQfDMMwMTJwqWHvH
nSSFznjFcfjTgVJTJjvT
RrBpcfSNpRBcFshrCtQrDGLPQb
GctcMldStGwPPbcLsQTV
jhnzDgnHnnfPVwHQTTLTds
gBgzDDhzvqdGGvrtdvtm
PPwRWVTvRvPVLwRpMlzmbmsbHWjbbs
dFTFCNfdjzjFjsMF
SgdffSTrnnqCgdqgcNrfSZqVvVLRDPJZQwJBLPQPtLZwJJ
HSzDQftHphTBHFhr
WMmJsMJNLWNPmmsncjMJcjtvwggVvhFFFhjrVwrjrjppvT
WWPsJJCWCtZZZRGC
SfFZQDRLgpLlRgQRRRFWTsbhBhgTjbWBjhshgw
tHVNGjtzzHvMMJJhrWJNrTrPbP
jvVVvGCGtCmjHtdHzQQfDpSSlDRnFfQFmR
ZBBPfVVPPrVmrWZJzNdPznbnbSzznP
gvgpGqFFMgMcGgwLwGpcJNZSTZbdbdzNNSlqbTnn
ssZLgsMLQvcFpVrhCsmtWHjrrW
nFvhRnWWzBRPHQqcqqCqmFbd
SJDJgVprLfDfbJmHmWWHQtJW
wVsgWVSTgLjfSsVjVBRvRNwwGGRhNZhRMZ
wTRrRBCTPTBPlgMqgHCqggHLgg
dmDzpQpBdmWmWzzzDFzjGNMSWLSgLVGGLVNSLgSg
JbFdmQQJpjpptQbdJJmDmdtnZhZRflTlnTrTlrhBwPPc
jwSwssQbwbStlhRgtsVstn
zPzFdFFZccPDgntzVHHgghRz
vVVdfFmZPDMWZZBmGQbwJJGCLwwMwJCS
PLLffLFqqNLwSffbnVzzRf
lsmgTggChrgDtZsZVblMVJMznVnwMBnb
hDZZTmZvhTgstvwNFdqpWQcqvP
mmWwpwqtmmHTqHpprRZQPPZLZWSFRFRB
gzcgscgbfvhRRNZQRRQPvr
JsjhcshCfJgrrpttTCTplH
TvNrvNrJfWWvtJLTHhvZZhQQwVGZZhbV
mFCPmBMPlPsPPBsMFPszbHQJwwHHhbZVQVzjhGHZ
qmlBsdCSmJBmsBBMnMMTDWDNcLrDprSNcWDNrW
HSnHCVqTddFHSVqFqdStSQGQwRzQCzCRzGRRGNPQwz
jhlBpgNvlfZjlfvmpgfgfBrMRRwMMPPLQMZRPMWMMZPQ
hflfpgjfBfDcpchlpvndbHFDdHqSqNVqNVdn
QBfhlVNfHSZHfVCVHSQfZfTCctdvdDTRtjDCtTRsjsvj
brrWWqzFWzwWbswDchhTRtDhjT
gbrFLnpzqrWPgqpLWrhnHGZlSfLfHfQBGVfHHNVG
hcFmVScmQBVhtcvfHLfvHSnbHRLn
lzQqlDqgpWPvJfRnlNJvww
PjCCCPgDgqGzmMtCQstZmBFs
GrnrHrmVVFMFhSSbSfhR
zjTqLtBjjdWdWTLshMZMDbPNRMSSqSPNfD
dTWjdtwWhjWTBzcnrpcwmQpCwcpw
BgtVBsgVVJhgGsSGJbghJqbsjLfZmLjmmtfZndNNZFFZNLdm
MHTlzlwHSvPvzMSPCTMQCNZdjjmnfQjdZfNLjF
HMDTwPDpzrTzpScBghDbVqRBgRBJsR
qjCsTmrrnnCmhcFrCjqmThRlbHGvJGvvvbRbbJRcQJRG
fVBBTfMdSZLNZgPdgglGRdbtbtJRllJWJtWb
TSZZSpNMPBLgpZLVgppBDFhzDmjrnzhjssFqDhzFnn
MDtDMWmMQmdzmMMqvlGfRcjzgpcPPjczPl
sZsHJWNJFJNbWFBhFBnnbVclbVVPglfcRggffGGpjl
srHWLNZZJdSLSdvDMw
BFqsPnsZcgnncggccqsqqpDPtDWPpPNTNSNThrWtSj
dQbfQQJJdwdmFRbJLRJdMrfDWrjpDrrprhDDDDNWWD
QGwdJmCFJZvGsvvncB
GRRNSjrffGTSPrNTffSgcJTwWJZbbZvwvwtVwWVJZv
CMcFsqmBQzMzshsBQBQvLWmJVwZtpJLJwJLvwt
qBzChnQlzMcRSnnRRjHGGg
dPPbPWNdTBbDpHPHpNsgzvFlglvHzvSFzCFF
fntqhGhRMhnnnGGCVMRhCVlSjswFvvjzjSzvQVSgsl
MMqqJJRCnMhZLCRhtPNrpZPDDWTDTrNBmD
VjVGNTNlNchVjNGRWrSWWtZtRrzncR
BTbbbwDmCDLTbDwfHmzZMmRrWtzPrZSZtrrz
qDLLqvwLBsfBbBdQFQTJhqThVTQJ
qBqPBGZflhrWznzZZdsnzv
FmHRcCCsCDwDbjtzdjWbdz
cTRwCCNHFNmwRgmFTNCFTJqqfqPJsPhPBlsrsfGf
JTNhhNrCTcWpJJcpWw
LLdLsfMsdStbtggLbnTpwfWDzpjvnnjzcz
MqSZsgbMbGFbZtgTSZSFSLhhmHlBQlmrmrPHHCPPZPBr
GgjjgpGvpJJtjgvPrJttsLjVwCsQsQNLsfLfMVCl
RddqZqzcZdSWcHdcqQfQLMwVflfNQNMQ
FncWTRcFlHWFmcFTgPJhGvDgpnnhhGtG
ZGPFJsPQCbZCCbgz
nrvrnGWTwwqTBRcpCRRg
DnGWDldNmSLSMdMQ
ScDmPPPmjmjjWgtdSmdmCnNqVQVVrNRTZTQTQDHHZGNH
wbbMhLvMJpRwJQrJGHqQGrQHqr
wzhspwpppswsMFLmtnzjdcmdPRWPtS
fPlLTtBlTjDbWcTMJcncWqss
LdvLrLpCRRQQmvhhVNpwRMWJZcFMZWJwJFFFwZsGcs
prrRVpvvRVNmRvvzHjHPlzPLSffbbPHS
bbCbzsQbBzbBFbFzFfJHfVJPfVPtzZttpm
hwvrjDDwjcDcvdnNvwdnwwDPpNtNMVtPpJpRmVMPfmtVZR
mmqDWjhcqhwvDDdTsSGGCbQGSBBLTGGl
wthtwrBQhhSrqJJVLMRPPtJLPF
vbmbZqfqgsfHmcllmLLJNLJMNFJNNvPGJn
jHlDscllClCCgbZzhWqBzBQzprWCqS
DwpDlwDwllhJwbDFNDwFPhDnQnZZzVVnnBrtNznzSSzZrr
MHPfWRTgWzTQmSTTZt
GdPjHLRMCfcGvwqDGbFq
NGdNpDPdNGLppLpwSNFFFDLwnnZbMnrHlHZQcHbcnHQcPcQM
jBgBsjssgjWJVGhBfhvJVnMlbnZQVrrrlHCcCZrl
sRBgvRffhhtvtBgvjgttWBJTLmFLLwdwzFdqLDDDqFGdLFTd
RBzRGRGBgnrPJrRGGWWbDggZTSTVZVFFSVZZdw
tfmLhjfshNNsqpppfjHdwVdQZbVbQFVVTmSVDV
jvqtwHqCNLNsNhNfjssLCzGJWcBRJRBCBBBrrWnGrG
cGDBdNFdNdDStNtGdGQGscDMvZjjZWMvMLjsjjLZjLsJLv
fbRnRzHVPClfRlbbmfHRlPvZvMLpqjvZJqZMpzMpZjBZ
gfbfBgCHCPVNhNGgTtFDcF
ljjvLZvvllFnlLJLJjLWFQrVssFpsGMFpNMGsr
qbHSqHTcHChqCmTSRqqBVpVdBQsVQQGHdMprWV
DbtcCqRhhTRmDnPMlLgfgztlZZ
QtMQzPZcbtGgTtLvtf
HwcDsdVVDnNSGLhwvffJvfTT
CSsdSHNmSDHcCqqcrrzQcC
nDNDfmMnmDJRNfJJdMDRBdwjcTtsSHvBTswTjwLtsCts
QcbQrZZggGWVVWbZZjStswCSCCttGCwLvC
ZgZgbzVglVbchVVrVhFWWpnnpNJmJMqDfMJnMRqNDDMD
WMfNvsjWGjsqFjqTRRQVJztDzVmJHbft
ZPhplcrLrmzQGzmz
ZPddLPlcSclhZChndMvgTjjWNGWMGWBj
nDLjMchghDcljfjffpfsqTmGCTGszGZVVzHV
wdFJPFrQRwSNbjVQCTGsHZHmHCHs
BRJJSddvdBrJwrBRNRFRSPRjvclLpWhpglgWpLplltnMgh
BbVRzMcpbMNMHMTJmWdljdlJjT
GtsqFnfvGSFqGfQvgnWWZlLlLjZWtWldPmlT
sSsFqsqsGghwQQmfGRHbbVczbwwBpBpHcw
BBFCBJCsGJBBgvgsvTlVhgNg
ZnLdjRQddLRnZrlScHRVTTSHhRvg
fnnjZLWdrnqdWrrPLddqVqBzGDJJFGCBDfJmbDzFMbmB`
