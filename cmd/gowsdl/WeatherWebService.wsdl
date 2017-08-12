<?xml version="1.0" encoding="utf-8"?>
<wsdl:definitions xmlns:soap="http://schemas.xmlsoap.org/wsdl/soap/" xmlns:tm="http://microsoft.com/wsdl/mime/textMatching/" xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:mime="http://schemas.xmlsoap.org/wsdl/mime/" xmlns:tns="http://WebXml.com.cn/" xmlns:s="http://www.w3.org/2001/XMLSchema" xmlns:soap12="http://schemas.xmlsoap.org/wsdl/soap12/" xmlns:http="http://schemas.xmlsoap.org/wsdl/http/" targetNamespace="http://WebXml.com.cn/" xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/">
  <wsdl:documentation xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/">&lt;a href="http://www.webxml.com.cn/" target="_blank"&gt;WebXml.com.cn&lt;/a&gt; &lt;strong&gt;天气预报 Web 服务，数据每2.5小时左右自动更新一次，准确可靠。包括 340 多个中国主要城市和 60 多个国外主要城市三日内的天气预报数据。&lt;/br&gt;此天气预报Web Services请不要用于任何商业目的，若有需要请&lt;a href="http://www.webxml.com.cn/zh_cn/contact_us.aspx" target="_blank"&gt;联系我们&lt;/a&gt;，欢迎技术交流。 QQ：8409035&lt;br /&gt;使用本站 WEB 服务请注明或链接本站：http://www.webxml.com.cn/ 感谢大家的支持&lt;/strong&gt;！&lt;br /&gt;&lt;span style="color:#999999;"&gt;通知：天气预报 WEB 服务如原来使用地址 http://www.onhap.com/WebServices/WeatherWebService.asmx 的，请改成现在使用的服务地址 http://www.webxml.com.cn/WebServices/WeatherWebService.asmx ，重新引用即可。&lt;/span&gt;&lt;br /&gt;&lt;br /&gt;&amp;nbsp;</wsdl:documentation>
  <wsdl:types>
    <s:schema elementFormDefault="qualified" targetNamespace="http://WebXml.com.cn/">
      <s:element name="getSupportCity">
        <s:complexType>
          <s:sequence>
            <s:element minOccurs="0" maxOccurs="1" name="byProvinceName" type="s:string" />
          </s:sequence>
        </s:complexType>
      </s:element>
      <s:element name="getSupportCityResponse">
        <s:complexType>
          <s:sequence>
            <s:element minOccurs="0" maxOccurs="1" name="getSupportCityResult" type="tns:ArrayOfString" />
          </s:sequence>
        </s:complexType>
      </s:element>
      <s:complexType name="ArrayOfString">
        <s:sequence>
          <s:element minOccurs="0" maxOccurs="unbounded" name="string" nillable="true" type="s:string" />
        </s:sequence>
      </s:complexType>
      <s:element name="getSupportProvince">
        <s:complexType />
      </s:element>
      <s:element name="getSupportProvinceResponse">
        <s:complexType>
          <s:sequence>
            <s:element minOccurs="0" maxOccurs="1" name="getSupportProvinceResult" type="tns:ArrayOfString" />
          </s:sequence>
        </s:complexType>
      </s:element>
      <s:element name="getSupportDataSet">
        <s:complexType />
      </s:element>
      <s:element name="getSupportDataSetResponse">
        <s:complexType>
          <s:sequence>
            <s:element minOccurs="0" maxOccurs="1" name="getSupportDataSetResult">
              <s:complexType>
                <s:sequence>
                  <s:element ref="s:schema" />
                  <s:any />
                </s:sequence>
              </s:complexType>
            </s:element>
          </s:sequence>
        </s:complexType>
      </s:element>
      <s:element name="getWeatherbyCityName">
        <s:complexType>
          <s:sequence>
            <s:element minOccurs="0" maxOccurs="1" name="theCityName" type="s:string" />
          </s:sequence>
        </s:complexType>
      </s:element>
      <s:element name="getWeatherbyCityNameResponse">
        <s:complexType>
          <s:sequence>
            <s:element minOccurs="0" maxOccurs="1" name="getWeatherbyCityNameResult" type="tns:ArrayOfString" />
          </s:sequence>
        </s:complexType>
      </s:element>
      <s:element name="getWeatherbyCityNamePro">
        <s:complexType>
          <s:sequence>
            <s:element minOccurs="0" maxOccurs="1" name="theCityName" type="s:string" />
            <s:element minOccurs="0" maxOccurs="1" name="theUserID" type="s:string" />
          </s:sequence>
        </s:complexType>
      </s:element>
      <s:element name="getWeatherbyCityNameProResponse">
        <s:complexType>
          <s:sequence>
            <s:element minOccurs="0" maxOccurs="1" name="getWeatherbyCityNameProResult" type="tns:ArrayOfString" />
          </s:sequence>
        </s:complexType>
      </s:element>
      <s:element name="ArrayOfString" nillable="true" type="tns:ArrayOfString" />
      <s:element name="DataSet" nillable="true">
        <s:complexType>
          <s:sequence>
            <s:element ref="s:schema" />
            <s:any />
          </s:sequence>
        </s:complexType>
      </s:element>
    </s:schema>
  </wsdl:types>
  <wsdl:message name="getSupportCitySoapIn">
    <wsdl:part name="parameters" element="tns:getSupportCity" />
  </wsdl:message>
  <wsdl:message name="getSupportCitySoapOut">
    <wsdl:part name="parameters" element="tns:getSupportCityResponse" />
  </wsdl:message>
  <wsdl:message name="getSupportProvinceSoapIn">
    <wsdl:part name="parameters" element="tns:getSupportProvince" />
  </wsdl:message>
  <wsdl:message name="getSupportProvinceSoapOut">
    <wsdl:part name="parameters" element="tns:getSupportProvinceResponse" />
  </wsdl:message>
  <wsdl:message name="getSupportDataSetSoapIn">
    <wsdl:part name="parameters" element="tns:getSupportDataSet" />
  </wsdl:message>
  <wsdl:message name="getSupportDataSetSoapOut">
    <wsdl:part name="parameters" element="tns:getSupportDataSetResponse" />
  </wsdl:message>
  <wsdl:message name="getWeatherbyCityNameSoapIn">
    <wsdl:part name="parameters" element="tns:getWeatherbyCityName" />
  </wsdl:message>
  <wsdl:message name="getWeatherbyCityNameSoapOut">
    <wsdl:part name="parameters" element="tns:getWeatherbyCityNameResponse" />
  </wsdl:message>
  <wsdl:message name="getWeatherbyCityNameProSoapIn">
    <wsdl:part name="parameters" element="tns:getWeatherbyCityNamePro" />
  </wsdl:message>
  <wsdl:message name="getWeatherbyCityNameProSoapOut">
    <wsdl:part name="parameters" element="tns:getWeatherbyCityNameProResponse" />
  </wsdl:message>
  <wsdl:message name="getSupportCityHttpGetIn">
    <wsdl:part name="byProvinceName" type="s:string" />
  </wsdl:message>
  <wsdl:message name="getSupportCityHttpGetOut">
    <wsdl:part name="Body" element="tns:ArrayOfString" />
  </wsdl:message>
  <wsdl:message name="getSupportProvinceHttpGetIn" />
  <wsdl:message name="getSupportProvinceHttpGetOut">
    <wsdl:part name="Body" element="tns:ArrayOfString" />
  </wsdl:message>
  <wsdl:message name="getSupportDataSetHttpGetIn" />
  <wsdl:message name="getSupportDataSetHttpGetOut">
    <wsdl:part name="Body" element="tns:DataSet" />
  </wsdl:message>
  <wsdl:message name="getWeatherbyCityNameHttpGetIn">
    <wsdl:part name="theCityName" type="s:string" />
  </wsdl:message>
  <wsdl:message name="getWeatherbyCityNameHttpGetOut">
    <wsdl:part name="Body" element="tns:ArrayOfString" />
  </wsdl:message>
  <wsdl:message name="getWeatherbyCityNameProHttpGetIn">
    <wsdl:part name="theCityName" type="s:string" />
    <wsdl:part name="theUserID" type="s:string" />
  </wsdl:message>
  <wsdl:message name="getWeatherbyCityNameProHttpGetOut">
    <wsdl:part name="Body" element="tns:ArrayOfString" />
  </wsdl:message>
  <wsdl:message name="getSupportCityHttpPostIn">
    <wsdl:part name="byProvinceName" type="s:string" />
  </wsdl:message>
  <wsdl:message name="getSupportCityHttpPostOut">
    <wsdl:part name="Body" element="tns:ArrayOfString" />
  </wsdl:message>
  <wsdl:message name="getSupportProvinceHttpPostIn" />
  <wsdl:message name="getSupportProvinceHttpPostOut">
    <wsdl:part name="Body" element="tns:ArrayOfString" />
  </wsdl:message>
  <wsdl:message name="getSupportDataSetHttpPostIn" />
  <wsdl:message name="getSupportDataSetHttpPostOut">
    <wsdl:part name="Body" element="tns:DataSet" />
  </wsdl:message>
  <wsdl:message name="getWeatherbyCityNameHttpPostIn">
    <wsdl:part name="theCityName" type="s:string" />
  </wsdl:message>
  <wsdl:message name="getWeatherbyCityNameHttpPostOut">
    <wsdl:part name="Body" element="tns:ArrayOfString" />
  </wsdl:message>
  <wsdl:message name="getWeatherbyCityNameProHttpPostIn">
    <wsdl:part name="theCityName" type="s:string" />
    <wsdl:part name="theUserID" type="s:string" />
  </wsdl:message>
  <wsdl:message name="getWeatherbyCityNameProHttpPostOut">
    <wsdl:part name="Body" element="tns:ArrayOfString" />
  </wsdl:message>
  <wsdl:portType name="WeatherWebServiceSoap">
    <wsdl:operation name="getSupportCity">
      <wsdl:documentation xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/">&lt;br /&gt;&lt;h3&gt;查询本天气预报Web Services支持的国内外城市或地区信息&lt;/h3&gt;&lt;p&gt;输入参数：byProvinceName = 指定的洲或国内的省份，若为ALL或空则表示返回全部城市；返回数据：一个一维字符串数组 String()，结构为：城市名称(城市代码)。&lt;/p&gt;&lt;br /&gt;</wsdl:documentation>
      <wsdl:input message="tns:getSupportCitySoapIn" />
      <wsdl:output message="tns:getSupportCitySoapOut" />
    </wsdl:operation>
    <wsdl:operation name="getSupportProvince">
      <wsdl:documentation xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/">&lt;br /&gt;&lt;h3&gt;获得本天气预报Web Services支持的洲、国内外省份和城市信息&lt;/h3&gt;&lt;p&gt;输入参数：无； 返回数据：一个一维字符串数组 String()，内容为洲或国内省份的名称。&lt;/p&gt;&lt;br /&gt;</wsdl:documentation>
      <wsdl:input message="tns:getSupportProvinceSoapIn" />
      <wsdl:output message="tns:getSupportProvinceSoapOut" />
    </wsdl:operation>
    <wsdl:operation name="getSupportDataSet">
      <wsdl:documentation xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/">&lt;br&gt;&lt;h3&gt;获得本天气预报Web Services支持的洲、国内外省份和城市信息&lt;/h3&gt;&lt;p&gt;输入参数：无；返回：DataSet 。DataSet.Tables(0) 为支持的洲和国内省份数据，DataSet.Tables(1) 为支持的国内外城市或地区数据。DataSet.Tables(0).Rows(i).Item("ID") 主键对应 DataSet.Tables(1).Rows(i).Item("ZoneID") 外键。&lt;br /&gt;Tables(0)：ID = ID主键，Zone = 支持的洲、省份；Tables(1)：ID 主键，ZoneID = 对应Tables(0)ID的外键，Area = 城市或地区，AreaCode = 城市或地区代码。&lt;/p&gt;&lt;br /&gt;</wsdl:documentation>
      <wsdl:input message="tns:getSupportDataSetSoapIn" />
      <wsdl:output message="tns:getSupportDataSetSoapOut" />
    </wsdl:operation>
    <wsdl:operation name="getWeatherbyCityName">
      <wsdl:documentation xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/">&lt;br&gt;&lt;h3&gt;根据城市或地区名称查询获得未来三天内天气情况、现在的天气实况、天气和生活指数&lt;/h3&gt;&lt;p&gt;调用方法如下：输入参数：theCityName = 城市中文名称(国外城市可用英文)或城市代码(不输入默认为上海市)，如：上海 或 58367，如有城市名称重复请使用城市代码查询(可通过 getSupportCity 或 getSupportDataSet 获得)；返回数据： 一个一维数组 String(22)，共有23个元素。&lt;br /&gt;String(0) 到 String(4)：省份，城市，城市代码，城市图片名称，最后更新时间。String(5) 到 String(11)：当天的 气温，概况，风向和风力，天气趋势开始图片名称(以下称：图标一)，天气趋势结束图片名称(以下称：图标二)，现在的天气实况，天气和生活指数。String(12) 到 String(16)：第二天的 气温，概况，风向和风力，图标一，图标二。String(17) 到 String(21)：第三天的 气温，概况，风向和风力，图标一，图标二。String(22) 被查询的城市或地区的介绍 &lt;br /&gt;&lt;a href="http://www.webxml.com.cn/images/weather.zip"&gt;下载天气图标&lt;img src="http://www.webxml.com.cn/images/download_w.gif" border="0" align="absbottom" /&gt;&lt;/a&gt;(包含大、中、小尺寸) &lt;a href="http://www.webxml.com.cn/zh_cn/weather_icon.aspx" target="_blank"&gt;天气图例说明&lt;/a&gt; &lt;a href="http://www.webxml.com.cn/files/weather_eg.zip"&gt;调用此天气预报Web Services实例下载&lt;/a&gt; (VB ASP.net 2.0)&lt;/p&gt;&lt;br /&gt;</wsdl:documentation>
      <wsdl:input message="tns:getWeatherbyCityNameSoapIn" />
      <wsdl:output message="tns:getWeatherbyCityNameSoapOut" />
    </wsdl:operation>
    <wsdl:operation name="getWeatherbyCityNamePro">
      <wsdl:documentation xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/">&lt;br&gt;&lt;h3&gt;根据城市或地区名称查询获得未来三天内天气情况、现在的天气实况、天气和生活指数（For商业用户）&lt;/h3&gt;&lt;p&gt;调用方法同 getWeatherbyCityName，输入参数：theUserID = 商业用户ID&lt;/p&gt;&lt;br /&gt;</wsdl:documentation>
      <wsdl:input message="tns:getWeatherbyCityNameProSoapIn" />
      <wsdl:output message="tns:getWeatherbyCityNameProSoapOut" />
    </wsdl:operation>
  </wsdl:portType>
  <wsdl:portType name="WeatherWebServiceHttpGet">
    <wsdl:operation name="getSupportCity">
      <wsdl:documentation xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/">&lt;br /&gt;&lt;h3&gt;查询本天气预报Web Services支持的国内外城市或地区信息&lt;/h3&gt;&lt;p&gt;输入参数：byProvinceName = 指定的洲或国内的省份，若为ALL或空则表示返回全部城市；返回数据：一个一维字符串数组 String()，结构为：城市名称(城市代码)。&lt;/p&gt;&lt;br /&gt;</wsdl:documentation>
      <wsdl:input message="tns:getSupportCityHttpGetIn" />
      <wsdl:output message="tns:getSupportCityHttpGetOut" />
    </wsdl:operation>
    <wsdl:operation name="getSupportProvince">
      <wsdl:documentation xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/">&lt;br /&gt;&lt;h3&gt;获得本天气预报Web Services支持的洲、国内外省份和城市信息&lt;/h3&gt;&lt;p&gt;输入参数：无； 返回数据：一个一维字符串数组 String()，内容为洲或国内省份的名称。&lt;/p&gt;&lt;br /&gt;</wsdl:documentation>
      <wsdl:input message="tns:getSupportProvinceHttpGetIn" />
      <wsdl:output message="tns:getSupportProvinceHttpGetOut" />
    </wsdl:operation>
    <wsdl:operation name="getSupportDataSet">
      <wsdl:documentation xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/">&lt;br&gt;&lt;h3&gt;获得本天气预报Web Services支持的洲、国内外省份和城市信息&lt;/h3&gt;&lt;p&gt;输入参数：无；返回：DataSet 。DataSet.Tables(0) 为支持的洲和国内省份数据，DataSet.Tables(1) 为支持的国内外城市或地区数据。DataSet.Tables(0).Rows(i).Item("ID") 主键对应 DataSet.Tables(1).Rows(i).Item("ZoneID") 外键。&lt;br /&gt;Tables(0)：ID = ID主键，Zone = 支持的洲、省份；Tables(1)：ID 主键，ZoneID = 对应Tables(0)ID的外键，Area = 城市或地区，AreaCode = 城市或地区代码。&lt;/p&gt;&lt;br /&gt;</wsdl:documentation>
      <wsdl:input message="tns:getSupportDataSetHttpGetIn" />
      <wsdl:output message="tns:getSupportDataSetHttpGetOut" />
    </wsdl:operation>
    <wsdl:operation name="getWeatherbyCityName">
      <wsdl:documentation xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/">&lt;br&gt;&lt;h3&gt;根据城市或地区名称查询获得未来三天内天气情况、现在的天气实况、天气和生活指数&lt;/h3&gt;&lt;p&gt;调用方法如下：输入参数：theCityName = 城市中文名称(国外城市可用英文)或城市代码(不输入默认为上海市)，如：上海 或 58367，如有城市名称重复请使用城市代码查询(可通过 getSupportCity 或 getSupportDataSet 获得)；返回数据： 一个一维数组 String(22)，共有23个元素。&lt;br /&gt;String(0) 到 String(4)：省份，城市，城市代码，城市图片名称，最后更新时间。String(5) 到 String(11)：当天的 气温，概况，风向和风力，天气趋势开始图片名称(以下称：图标一)，天气趋势结束图片名称(以下称：图标二)，现在的天气实况，天气和生活指数。String(12) 到 String(16)：第二天的 气温，概况，风向和风力，图标一，图标二。String(17) 到 String(21)：第三天的 气温，概况，风向和风力，图标一，图标二。String(22) 被查询的城市或地区的介绍 &lt;br /&gt;&lt;a href="http://www.webxml.com.cn/images/weather.zip"&gt;下载天气图标&lt;img src="http://www.webxml.com.cn/images/download_w.gif" border="0" align="absbottom" /&gt;&lt;/a&gt;(包含大、中、小尺寸) &lt;a href="http://www.webxml.com.cn/zh_cn/weather_icon.aspx" target="_blank"&gt;天气图例说明&lt;/a&gt; &lt;a href="http://www.webxml.com.cn/files/weather_eg.zip"&gt;调用此天气预报Web Services实例下载&lt;/a&gt; (VB ASP.net 2.0)&lt;/p&gt;&lt;br /&gt;</wsdl:documentation>
      <wsdl:input message="tns:getWeatherbyCityNameHttpGetIn" />
      <wsdl:output message="tns:getWeatherbyCityNameHttpGetOut" />
    </wsdl:operation>
    <wsdl:operation name="getWeatherbyCityNamePro">
      <wsdl:documentation xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/">&lt;br&gt;&lt;h3&gt;根据城市或地区名称查询获得未来三天内天气情况、现在的天气实况、天气和生活指数（For商业用户）&lt;/h3&gt;&lt;p&gt;调用方法同 getWeatherbyCityName，输入参数：theUserID = 商业用户ID&lt;/p&gt;&lt;br /&gt;</wsdl:documentation>
      <wsdl:input message="tns:getWeatherbyCityNameProHttpGetIn" />
      <wsdl:output message="tns:getWeatherbyCityNameProHttpGetOut" />
    </wsdl:operation>
  </wsdl:portType>
  <wsdl:portType name="WeatherWebServiceHttpPost">
    <wsdl:operation name="getSupportCity">
      <wsdl:documentation xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/">&lt;br /&gt;&lt;h3&gt;查询本天气预报Web Services支持的国内外城市或地区信息&lt;/h3&gt;&lt;p&gt;输入参数：byProvinceName = 指定的洲或国内的省份，若为ALL或空则表示返回全部城市；返回数据：一个一维字符串数组 String()，结构为：城市名称(城市代码)。&lt;/p&gt;&lt;br /&gt;</wsdl:documentation>
      <wsdl:input message="tns:getSupportCityHttpPostIn" />
      <wsdl:output message="tns:getSupportCityHttpPostOut" />
    </wsdl:operation>
    <wsdl:operation name="getSupportProvince">
      <wsdl:documentation xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/">&lt;br /&gt;&lt;h3&gt;获得本天气预报Web Services支持的洲、国内外省份和城市信息&lt;/h3&gt;&lt;p&gt;输入参数：无； 返回数据：一个一维字符串数组 String()，内容为洲或国内省份的名称。&lt;/p&gt;&lt;br /&gt;</wsdl:documentation>
      <wsdl:input message="tns:getSupportProvinceHttpPostIn" />
      <wsdl:output message="tns:getSupportProvinceHttpPostOut" />
    </wsdl:operation>
    <wsdl:operation name="getSupportDataSet">
      <wsdl:documentation xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/">&lt;br&gt;&lt;h3&gt;获得本天气预报Web Services支持的洲、国内外省份和城市信息&lt;/h3&gt;&lt;p&gt;输入参数：无；返回：DataSet 。DataSet.Tables(0) 为支持的洲和国内省份数据，DataSet.Tables(1) 为支持的国内外城市或地区数据。DataSet.Tables(0).Rows(i).Item("ID") 主键对应 DataSet.Tables(1).Rows(i).Item("ZoneID") 外键。&lt;br /&gt;Tables(0)：ID = ID主键，Zone = 支持的洲、省份；Tables(1)：ID 主键，ZoneID = 对应Tables(0)ID的外键，Area = 城市或地区，AreaCode = 城市或地区代码。&lt;/p&gt;&lt;br /&gt;</wsdl:documentation>
      <wsdl:input message="tns:getSupportDataSetHttpPostIn" />
      <wsdl:output message="tns:getSupportDataSetHttpPostOut" />
    </wsdl:operation>
    <wsdl:operation name="getWeatherbyCityName">
      <wsdl:documentation xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/">&lt;br&gt;&lt;h3&gt;根据城市或地区名称查询获得未来三天内天气情况、现在的天气实况、天气和生活指数&lt;/h3&gt;&lt;p&gt;调用方法如下：输入参数：theCityName = 城市中文名称(国外城市可用英文)或城市代码(不输入默认为上海市)，如：上海 或 58367，如有城市名称重复请使用城市代码查询(可通过 getSupportCity 或 getSupportDataSet 获得)；返回数据： 一个一维数组 String(22)，共有23个元素。&lt;br /&gt;String(0) 到 String(4)：省份，城市，城市代码，城市图片名称，最后更新时间。String(5) 到 String(11)：当天的 气温，概况，风向和风力，天气趋势开始图片名称(以下称：图标一)，天气趋势结束图片名称(以下称：图标二)，现在的天气实况，天气和生活指数。String(12) 到 String(16)：第二天的 气温，概况，风向和风力，图标一，图标二。String(17) 到 String(21)：第三天的 气温，概况，风向和风力，图标一，图标二。String(22) 被查询的城市或地区的介绍 &lt;br /&gt;&lt;a href="http://www.webxml.com.cn/images/weather.zip"&gt;下载天气图标&lt;img src="http://www.webxml.com.cn/images/download_w.gif" border="0" align="absbottom" /&gt;&lt;/a&gt;(包含大、中、小尺寸) &lt;a href="http://www.webxml.com.cn/zh_cn/weather_icon.aspx" target="_blank"&gt;天气图例说明&lt;/a&gt; &lt;a href="http://www.webxml.com.cn/files/weather_eg.zip"&gt;调用此天气预报Web Services实例下载&lt;/a&gt; (VB ASP.net 2.0)&lt;/p&gt;&lt;br /&gt;</wsdl:documentation>
      <wsdl:input message="tns:getWeatherbyCityNameHttpPostIn" />
      <wsdl:output message="tns:getWeatherbyCityNameHttpPostOut" />
    </wsdl:operation>
    <wsdl:operation name="getWeatherbyCityNamePro">
      <wsdl:documentation xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/">&lt;br&gt;&lt;h3&gt;根据城市或地区名称查询获得未来三天内天气情况、现在的天气实况、天气和生活指数（For商业用户）&lt;/h3&gt;&lt;p&gt;调用方法同 getWeatherbyCityName，输入参数：theUserID = 商业用户ID&lt;/p&gt;&lt;br /&gt;</wsdl:documentation>
      <wsdl:input message="tns:getWeatherbyCityNameProHttpPostIn" />
      <wsdl:output message="tns:getWeatherbyCityNameProHttpPostOut" />
    </wsdl:operation>
  </wsdl:portType>
  <wsdl:binding name="WeatherWebServiceSoap" type="tns:WeatherWebServiceSoap">
    <soap:binding transport="http://schemas.xmlsoap.org/soap/http" />
    <wsdl:operation name="getSupportCity">
      <soap:operation soapAction="http://WebXml.com.cn/getSupportCity" style="document" />
      <wsdl:input>
        <soap:body use="literal" />
      </wsdl:input>
      <wsdl:output>
        <soap:body use="literal" />
      </wsdl:output>
    </wsdl:operation>
    <wsdl:operation name="getSupportProvince">
      <soap:operation soapAction="http://WebXml.com.cn/getSupportProvince" style="document" />
      <wsdl:input>
        <soap:body use="literal" />
      </wsdl:input>
      <wsdl:output>
        <soap:body use="literal" />
      </wsdl:output>
    </wsdl:operation>
    <wsdl:operation name="getSupportDataSet">
      <soap:operation soapAction="http://WebXml.com.cn/getSupportDataSet" style="document" />
      <wsdl:input>
        <soap:body use="literal" />
      </wsdl:input>
      <wsdl:output>
        <soap:body use="literal" />
      </wsdl:output>
    </wsdl:operation>
    <wsdl:operation name="getWeatherbyCityName">
      <soap:operation soapAction="http://WebXml.com.cn/getWeatherbyCityName" style="document" />
      <wsdl:input>
        <soap:body use="literal" />
      </wsdl:input>
      <wsdl:output>
        <soap:body use="literal" />
      </wsdl:output>
    </wsdl:operation>
    <wsdl:operation name="getWeatherbyCityNamePro">
      <soap:operation soapAction="http://WebXml.com.cn/getWeatherbyCityNamePro" style="document" />
      <wsdl:input>
        <soap:body use="literal" />
      </wsdl:input>
      <wsdl:output>
        <soap:body use="literal" />
      </wsdl:output>
    </wsdl:operation>
  </wsdl:binding>
  <wsdl:binding name="WeatherWebServiceSoap12" type="tns:WeatherWebServiceSoap">
    <soap12:binding transport="http://schemas.xmlsoap.org/soap/http" />
    <wsdl:operation name="getSupportCity">
      <soap12:operation soapAction="http://WebXml.com.cn/getSupportCity" style="document" />
      <wsdl:input>
        <soap12:body use="literal" />
      </wsdl:input>
      <wsdl:output>
        <soap12:body use="literal" />
      </wsdl:output>
    </wsdl:operation>
    <wsdl:operation name="getSupportProvince">
      <soap12:operation soapAction="http://WebXml.com.cn/getSupportProvince" style="document" />
      <wsdl:input>
        <soap12:body use="literal" />
      </wsdl:input>
      <wsdl:output>
        <soap12:body use="literal" />
      </wsdl:output>
    </wsdl:operation>
    <wsdl:operation name="getSupportDataSet">
      <soap12:operation soapAction="http://WebXml.com.cn/getSupportDataSet" style="document" />
      <wsdl:input>
        <soap12:body use="literal" />
      </wsdl:input>
      <wsdl:output>
        <soap12:body use="literal" />
      </wsdl:output>
    </wsdl:operation>
    <wsdl:operation name="getWeatherbyCityName">
      <soap12:operation soapAction="http://WebXml.com.cn/getWeatherbyCityName" style="document" />
      <wsdl:input>
        <soap12:body use="literal" />
      </wsdl:input>
      <wsdl:output>
        <soap12:body use="literal" />
      </wsdl:output>
    </wsdl:operation>
    <wsdl:operation name="getWeatherbyCityNamePro">
      <soap12:operation soapAction="http://WebXml.com.cn/getWeatherbyCityNamePro" style="document" />
      <wsdl:input>
        <soap12:body use="literal" />
      </wsdl:input>
      <wsdl:output>
        <soap12:body use="literal" />
      </wsdl:output>
    </wsdl:operation>
  </wsdl:binding>
  <wsdl:binding name="WeatherWebServiceHttpGet" type="tns:WeatherWebServiceHttpGet">
    <http:binding verb="GET" />
    <wsdl:operation name="getSupportCity">
      <http:operation location="/getSupportCity" />
      <wsdl:input>
        <http:urlEncoded />
      </wsdl:input>
      <wsdl:output>
        <mime:mimeXml part="Body" />
      </wsdl:output>
    </wsdl:operation>
    <wsdl:operation name="getSupportProvince">
      <http:operation location="/getSupportProvince" />
      <wsdl:input>
        <http:urlEncoded />
      </wsdl:input>
      <wsdl:output>
        <mime:mimeXml part="Body" />
      </wsdl:output>
    </wsdl:operation>
    <wsdl:operation name="getSupportDataSet">
      <http:operation location="/getSupportDataSet" />
      <wsdl:input>
        <http:urlEncoded />
      </wsdl:input>
      <wsdl:output>
        <mime:mimeXml part="Body" />
      </wsdl:output>
    </wsdl:operation>
    <wsdl:operation name="getWeatherbyCityName">
      <http:operation location="/getWeatherbyCityName" />
      <wsdl:input>
        <http:urlEncoded />
      </wsdl:input>
      <wsdl:output>
        <mime:mimeXml part="Body" />
      </wsdl:output>
    </wsdl:operation>
    <wsdl:operation name="getWeatherbyCityNamePro">
      <http:operation location="/getWeatherbyCityNamePro" />
      <wsdl:input>
        <http:urlEncoded />
      </wsdl:input>
      <wsdl:output>
        <mime:mimeXml part="Body" />
      </wsdl:output>
    </wsdl:operation>
  </wsdl:binding>
  <wsdl:binding name="WeatherWebServiceHttpPost" type="tns:WeatherWebServiceHttpPost">
    <http:binding verb="POST" />
    <wsdl:operation name="getSupportCity">
      <http:operation location="/getSupportCity" />
      <wsdl:input>
        <mime:content type="application/x-www-form-urlencoded" />
      </wsdl:input>
      <wsdl:output>
        <mime:mimeXml part="Body" />
      </wsdl:output>
    </wsdl:operation>
    <wsdl:operation name="getSupportProvince">
      <http:operation location="/getSupportProvince" />
      <wsdl:input>
        <mime:content type="application/x-www-form-urlencoded" />
      </wsdl:input>
      <wsdl:output>
        <mime:mimeXml part="Body" />
      </wsdl:output>
    </wsdl:operation>
    <wsdl:operation name="getSupportDataSet">
      <http:operation location="/getSupportDataSet" />
      <wsdl:input>
        <mime:content type="application/x-www-form-urlencoded" />
      </wsdl:input>
      <wsdl:output>
        <mime:mimeXml part="Body" />
      </wsdl:output>
    </wsdl:operation>
    <wsdl:operation name="getWeatherbyCityName">
      <http:operation location="/getWeatherbyCityName" />
      <wsdl:input>
        <mime:content type="application/x-www-form-urlencoded" />
      </wsdl:input>
      <wsdl:output>
        <mime:mimeXml part="Body" />
      </wsdl:output>
    </wsdl:operation>
    <wsdl:operation name="getWeatherbyCityNamePro">
      <http:operation location="/getWeatherbyCityNamePro" />
      <wsdl:input>
        <mime:content type="application/x-www-form-urlencoded" />
      </wsdl:input>
      <wsdl:output>
        <mime:mimeXml part="Body" />
      </wsdl:output>
    </wsdl:operation>
  </wsdl:binding>
  <wsdl:service name="WeatherWebService">
    <wsdl:documentation xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/">&lt;a href="http://www.webxml.com.cn/" target="_blank"&gt;WebXml.com.cn&lt;/a&gt; &lt;strong&gt;天气预报 Web 服务，数据每2.5小时左右自动更新一次，准确可靠。包括 340 多个中国主要城市和 60 多个国外主要城市三日内的天气预报数据。&lt;/br&gt;此天气预报Web Services请不要用于任何商业目的，若有需要请&lt;a href="http://www.webxml.com.cn/zh_cn/contact_us.aspx" target="_blank"&gt;联系我们&lt;/a&gt;，欢迎技术交流。 QQ：8409035&lt;br /&gt;使用本站 WEB 服务请注明或链接本站：http://www.webxml.com.cn/ 感谢大家的支持&lt;/strong&gt;！&lt;br /&gt;&lt;span style="color:#999999;"&gt;通知：天气预报 WEB 服务如原来使用地址 http://www.onhap.com/WebServices/WeatherWebService.asmx 的，请改成现在使用的服务地址 http://www.webxml.com.cn/WebServices/WeatherWebService.asmx ，重新引用即可。&lt;/span&gt;&lt;br /&gt;&lt;br /&gt;&amp;nbsp;</wsdl:documentation>
    <wsdl:port name="WeatherWebServiceSoap" binding="tns:WeatherWebServiceSoap">
      <soap:address location="http://www.webxml.com.cn/WebServices/WeatherWebService.asmx" />
    </wsdl:port>
    <wsdl:port name="WeatherWebServiceSoap12" binding="tns:WeatherWebServiceSoap12">
      <soap12:address location="http://www.webxml.com.cn/WebServices/WeatherWebService.asmx" />
    </wsdl:port>
    <wsdl:port name="WeatherWebServiceHttpGet" binding="tns:WeatherWebServiceHttpGet">
      <http:address location="http://www.webxml.com.cn/WebServices/WeatherWebService.asmx" />
    </wsdl:port>
    <wsdl:port name="WeatherWebServiceHttpPost" binding="tns:WeatherWebServiceHttpPost">
      <http:address location="http://www.webxml.com.cn/WebServices/WeatherWebService.asmx" />
    </wsdl:port>
  </wsdl:service>
</wsdl:definitions>