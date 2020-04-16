package athenadriver

// getPriceOneByte to get the USD price per 1 Byte
// >>> 5.0/ (1024**4)
// 4.547473508864641e-12
// https://calculator.aws/pricing/2.0/meteredUnitMaps/athena/USD/current/athena.json
/**
{
   "manifest":{
      "serviceId":"athena",
      "accessType":"publish",
      "esIndex":"plc-athena-usd-20200313210020",
      "hawkFilePublicationDate":"2019-12-13T23:09:03Z",
      "currencyCode":"USD",
      "source":"athena"
   },
   "sets":{

   },
   "regions":{
      "AWS GovCloud (US)":{
         "Athena Queries":{
            "rateCode":"KHE6GFM67M54TKPD.JRTCKXETXF.6YS6EN2CT7",
            "price":"5.0000000000"
         }
      },
      "AWS GovCloud (US-East)":{
         "Athena Queries":{
            "rateCode":"ZTWT3PF4VA2UGXFT.JRTCKXETXF.6YS6EN2CT7",
            "price":"5.0000000000"
         }
      },
      "Asia Pacific (Hong Kong)":{
         "Athena Queries":{
            "rateCode":"CPTQKGPZUW93BMS3.JRTCKXETXF.6YS6EN2CT7",
            "price":"5.5000000000"
         }
      },
      "Asia Pacific (Mumbai)":{
         "Athena Queries":{
            "rateCode":"RQN88TYRT35JXK3M.JRTCKXETXF.6YS6EN2CT7",
            "price":"5.0000000000"
         }
      },
      "Asia Pacific (Seoul)":{
         "Athena Queries":{
            "rateCode":"CDJE83QEK4W8A85N.JRTCKXETXF.6YS6EN2CT7",
            "price":"5.0000000000"
         }
      },
      "Asia Pacific (Singapore)":{
         "Athena Queries":{
            "rateCode":"6YV9CGVD72AT3N83.JRTCKXETXF.6YS6EN2CT7",
            "price":"5.0000000000"
         }
      },
      "Asia Pacific (Sydney)":{
         "Athena Queries":{
            "rateCode":"55T4EWW7D53ERNFG.JRTCKXETXF.6YS6EN2CT7",
            "price":"5.0000000000"
         }
      },
      "Asia Pacific (Tokyo)":{
         "Athena Queries":{
            "rateCode":"6R553MKD6KWECHEK.JRTCKXETXF.6YS6EN2CT7",
            "price":"5.0000000000"
         }
      },
      "Canada (Central)":{
         "Athena Queries":{
            "rateCode":"ZR4NYZYR9SSMRTUV.JRTCKXETXF.6YS6EN2CT7",
            "price":"5.5000000000"
         }
      },
      "EU (Frankfurt)":{
         "Athena Queries":{
            "rateCode":"VFMVTH8MZDQM2MKA.JRTCKXETXF.6YS6EN2CT7",
            "price":"5.0000000000"
         }
      },
      "EU (Ireland)":{
         "Athena Queries":{
            "rateCode":"ZSVWQKQ6RCNFKBU6.JRTCKXETXF.6YS6EN2CT7",
            "price":"5.0000000000"
         }
      },
      "EU (London)":{
         "Athena Queries":{
            "rateCode":"FX2HWFNAPUT65ZJR.JRTCKXETXF.6YS6EN2CT7",
            "price":"5.0000000000"
         }
      },
      "EU (Paris)":{
         "Athena Queries":{
            "rateCode":"QJG3MDDH8VGVDPFU.JRTCKXETXF.6YS6EN2CT7",
            "price":"7.0000000000"
         }
      },
      "EU (Stockholm)":{
         "Athena Queries":{
            "rateCode":"UJ5FAC3WF3PAZKV9.JRTCKXETXF.6YS6EN2CT7",
            "price":"5.0000000000"
         }
      },
      "Middle East (Bahrain)":{
         "Athena Queries":{
            "rateCode":"9YQEMAJXRGSSVJPZ.JRTCKXETXF.6YS6EN2CT7",
            "price":"6.5000000000"
         }
      },
      "South America (Sao Paulo)":{
         "Athena Queries":{
            "rateCode":"7WNS6XGWHAMEK4ZQ.JRTCKXETXF.6YS6EN2CT7",
            "price":"9.0000000000"
         }
      },
      "US East (N. Virginia)":{
         "Athena Queries":{
            "rateCode":"B6WHE3FUNCDQVVX4.JRTCKXETXF.6YS6EN2CT7",
            "price":"5.0000000000"
         }
      },
      "US East (Ohio)":{
         "Athena Queries":{
            "rateCode":"E24RR3GXPVGYBMFP.JRTCKXETXF.6YS6EN2CT7",
            "price":"5.0000000000"
         }
      },
      "US West (N. California)":{
         "Athena Queries":{
            "rateCode":"BE4MTCX7QUSVD6PE.JRTCKXETXF.6YS6EN2CT7",
            "price":"6.7500000000"
         }
      },
      "US West (Oregon)":{
         "Athena Queries":{
            "rateCode":"WDJ89D3QKGNDVVCY.JRTCKXETXF.6YS6EN2CT7",
            "price":"5.0000000000"
         }
      }
   }
}
*/
func getPriceOneByte() float64 {
	return 4.547473508864641e-12
}

func getPrice10MB() float64 {
	return 10 * 1024 * 1024 * getPriceOneByte()
}
