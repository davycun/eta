package ctype_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
	"github.com/twpayne/go-geom/encoding/wkbhex"
	"testing"
)

func TestGeometry(t *testing.T) {

	var (
		s1 = `0106000000010000000103000000010000002B00000028EE75483C615E40409966AB20313F4038FA2C2A3D615E4090B82B6818313F40300C301744615E40E0B55C211F313F40B8E503D243615E40206E0DE122313F40D8A0128B43615E401081E11826313F40206A103343615E4020CE83DE2A313F406036A6BF4A615E40607D450D36313F40D8A669FB4B615E40600CF86042313F40180C25664F615E40103D9E353E313F40B8FA6A3E3E615E40D06487D592303F40C0E1258C3D615E4030843AC68B303F40B8B617113A615E405071EEB28F303F408863952F2A615E40A0D7F7EE97303F4038E9F4DD28615E4080D8476BA3303F40283B29FD1C615E40A0E6B7BA09313F402056FE1A19615E40E0030D7B29313F40B07D6EA216615E40C0BCDDAE3D313F400819331C12615E40107D1A3363313F4080015EB20D615E40101EEF758A313F40D01924A31B615E401071E0C1CD313F4020D0A99126615E40C0717F8702323F400807B30733615E4080900F193D323F40D8E6F61543615E40A0EFDC3C08323F4068E4D5342D615E4010901893B7313F404040DBF61F615E4000FB464BC8313F40C0ECCD2B1B615E4050D8B6FDB2313F40C02462C31A615E40807D1FED83313F40D8D02D5112615E40E0079C7477313F406895B9B01A615E40E08875E12C313F4060F391B120615E4070409319FA303F4088F0088C22615E4080342E6CEA303F40B882AEC129615E4010A01070AD303F40B0B1BD9131615E400028805AB8303F4068A3361B32615E404097807EB0303F40388B592735615E40A09C5C72AD303F40D01A619E36615E40F032B8BABB303F40E0AC71BC3F615E40703D7DFEC7303F409829EBD941615E4030D2AFE7CA303F40903FD71240615E4050ED1BC5DD303F4068CA87AF3D615E4060FB78C1DA303F40D082A1D337615E4020A5CA5D10313F40786569E636615E40B00386C518313F4028EE75483C615E40409966AB20313F40`
		//s2 = `0106000020E6100000010000000103000000010000002B00000028EE75483C615E40409966AB20313F4038FA2C2A3D615E4090B82B6818313F40300C301744615E40E0B55C211F313F40B8E503D243615E40206E0DE122313F40D8A0128B43615E401081E11826313F40206A103343615E4020CE83DE2A313F406036A6BF4A615E40607D450D36313F40D8A669FB4B615E40600CF86042313F40180C25664F615E40103D9E353E313F40B8FA6A3E3E615E40D06487D592303F40C0E1258C3D615E4030843AC68B303F40B8B617113A615E405071EEB28F303F408863952F2A615E40A0D7F7EE97303F4038E9F4DD28615E4080D8476BA3303F40283B29FD1C615E40A0E6B7BA09313F402056FE1A19615E40E0030D7B29313F40B07D6EA216615E40C0BCDDAE3D313F400819331C12615E40107D1A3363313F4080015EB20D615E40101EEF758A313F40D01924A31B615E401071E0C1CD313F4020D0A99126615E40C0717F8702323F400807B30733615E4080900F193D323F40D8E6F61543615E40A0EFDC3C08323F4068E4D5342D615E4010901893B7313F404040DBF61F615E4000FB464BC8313F40C0ECCD2B1B615E4050D8B6FDB2313F40C02462C31A615E40807D1FED83313F40D8D02D5112615E40E0079C7477313F406895B9B01A615E40E08875E12C313F4060F391B120615E4070409319FA303F4088F0088C22615E4080342E6CEA303F40B882AEC129615E4010A01070AD303F40B0B1BD9131615E400028805AB8303F4068A3361B32615E404097807EB0303F40388B592735615E40A09C5C72AD303F40D01A619E36615E40F032B8BABB303F40E0AC71BC3F615E40703D7DFEC7303F409829EBD941615E4030D2AFE7CA303F40903FD71240615E4050ED1BC5DD303F4068CA87AF3D615E4060FB78C1DA303F40D082A1D337615E4020A5CA5D10313F40786569E636615E40B00386C518313F4028EE75483C615E40409966AB20313F40`
	)

	gm1, err1 := ewkbhex.Decode(s1)
	gm2, err2 := wkbhex.Decode(s1)

	assert.Equal(t, len(gm1.Ends()), len(gm2.Ends()))
	assert.Nil(t, err1)
	assert.Nil(t, err2)

}
