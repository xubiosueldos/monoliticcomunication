package monoliticComunication

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/xubiosueldos/conexionBD/Concepto/structConcepto"

	"github.com/xubiosueldos/conexionBD/Autenticacion/structAutenticacion"
	"github.com/xubiosueldos/conexionBD/Helper/structHelper"
	"github.com/xubiosueldos/conexionBD/Legajo/structLegajo"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
	"github.com/xubiosueldos/framework"
	"github.com/xubiosueldos/framework/configuracion"
)

type requestMono struct {
	Value interface{}
	Error error
}

type strCuentaImporte struct {
	Cuentaid      int     `json:"cuentaid"`
	Importecuenta float32 `json:"importecuenta"`
}

type strLiquidacionContabilizar struct {
	requestMonolitico          strRequestMonolitico
	Descripcion                string             `json:"descripcion"`
	Cuentasimportes            []strCuentaImporte `json:"cuentasimportes"`
	Asientomanualtransaccionid int                `json:"asientomanualtransaccionid"`
}

type StrDatosAsientoContableManual struct {
	Asientocontablemanualid     int    `json:"asientocontablemanualid"`
	Asientocontablemanualnombre string `json:"asientocontablemanualnombre"`
}

func conectarconMonolitico(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, view string, columnid string, id string, options string) string {
	strReqMonolitico := llenarstructRequestMonolitico(tokenAutenticacion, id, options, view, columnid)
	str := reqMonolitico(w, r, view, columnid, strReqMonolitico)

	return str

}

func llenarstructRequestMonolitico(tokenAutenticacion *structAutenticacion.Security, id string, options string, view string, columnid string) *strRequestMonolitico {
	var strReqMonolitico strRequestMonolitico
	token := *tokenAutenticacion
	strReqMonolitico.Options = options
	strReqMonolitico.Tenant = token.Tenant
	strReqMonolitico.Token = token.Token
	strReqMonolitico.Username = token.Username
	strReqMonolitico.Id = id
	strReqMonolitico.View = view
	strReqMonolitico.Columnid = columnid

	return &strReqMonolitico
}

func reqMonolitico(w http.ResponseWriter, r *http.Request, view string, columnid string, structDinamico interface{}) string {

	url := configuracion.GetUrlMonolitico() + "MonoliticComunicationGoServlet"

	pagesJson, err := json.Marshal(structDinamico)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	fmt.Println("URL:>", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(pagesJson))

	if err != nil {
		fmt.Println("Error: ", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	str := string(body)
	fmt.Println("BYTES RECIBIDOS :", len(str))

	return str
}

func Obtenercentrodecosto(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, id string) *structLegajo.Centrodecosto {
	var centroDeCosto structLegajo.Centrodecosto
	str := conectarconMonolitico(w, r, tokenAutenticacion, "nxvcentrodecosto", "centrodecostoid", id, "CANQUERY")
	json.Unmarshal([]byte(str), &centroDeCosto)
	return &centroDeCosto
}

func Checkexistecentrodecosto(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, id string) *requestMono {
	var s requestMono
	centroDeCosto := Obtenercentrodecosto(w, r, tokenAutenticacion, id)
	if centroDeCosto.ID == 0 {
		s.Error = errors.New("El Centro de costo con ID: " + id + " no existe")
	}

	return &s
}

func Obtenercuenta(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, id string) *structConcepto.Cuenta {
	var cuenta structConcepto.Cuenta
	str := conectarconMonolitico(w, r, tokenAutenticacion, "nxvcuenta", "cuentaid", id, "CANQUERY")
	json.Unmarshal([]byte(str), &cuenta)
	return &cuenta
}

func Checkexistecuenta(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, id string) *requestMono {
	var s requestMono

	cuenta := Obtenercuenta(w, r, tokenAutenticacion, id)
	if cuenta.ID == 0 {
		s.Error = errors.New("La Cuenta con ID: " + id + " no existe")
	}

	return &s
}

func Obtenerbanco(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, id string) *structLiquidacion.Banco {
	var banco structLiquidacion.Banco
	str := conectarconMonolitico(w, r, tokenAutenticacion, "nxvbanco", "bancoid", id, "CANQUERY")
	json.Unmarshal([]byte(str), &banco)
	return &banco
}

func Checkexistebanco(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, id string) *requestMono {
	var s requestMono
	banco := Obtenerbanco(w, r, tokenAutenticacion, id)
	if banco.ID == 0 {
		s.Error = errors.New("El Banco con ID: " + id + " no existe")
	}

	return &s
}

func Gethelpers(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, codigo string, id string) *requestMono {

	var s requestMono
	view := "nxv" + codigo
	columnid := codigo + "id"
	str := conectarconMonolitico(w, r, tokenAutenticacion, view, columnid, id, "HLP")

	var dataHelper []structHelper.Helper
	json.Unmarshal([]byte(str), &dataHelper)
	framework.RespondJSON(w, http.StatusOK, dataHelper)

	return &s
}

func Obtenerdatosempresa(w http.ResponseWriter, r *http.Request, tokenAutenticacion *structAutenticacion.Security, codigo string, id string) *requestMono {
	var emp requestMono
	str := conectarconMonolitico(w, r, tokenAutenticacion, "fafempresa", "empresaid", id, "")

	var dataEmpresa structHelper.Empresa
	json.Unmarshal([]byte(str), &dataEmpresa)

	framework.RespondJSON(w, http.StatusOK, dataEmpresa)

	return &emp

}

func CheckAuthenticationMonolitico(tokenEncode string, r *http.Request) bool {

	infoUserValida := false
	var prueba []byte = []byte("xubiosueldosimplementadocongo")
	tokenSecurity := base64.StdEncoding.EncodeToString(prueba)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	url := configuracion.GetUrlMonolitico() + "SecurityAuthenticationGo"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Add("Authorization", tokenEncode)
	req.Header.Add("SecurityToken", tokenSecurity)

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusAccepted {
		infoUserValida = true
	}

	return infoUserValida
}

func requestMonoliticoContabilizarLiquidaciones(w http.ResponseWriter, r *http.Request, cuentasImportes []strCuentaImporte, tokenAutenticacion *structAutenticacion.Security, descripcion string, id string, options string, codigo string) string {

	var strLiquidacionContabilizar strLiquidacionContabilizar
	strReqMonolitico := llenarstructRequestMonolitico(tokenAutenticacion, "", "", id, options)

	if descripcion == "" {
		descripcion = framework.Descripcionasientomanualcontableliquidacionescontabilizadas
	}
	strLiquidacionContabilizar.requestMonolitico = *strReqMonolitico
	strLiquidacionContabilizar.Descripcion = descripcion
	strLiquidacionContabilizar.Cuentasimportes = cuentasImportes

	str := reqMonolitico(w, r, codigo, "", strLiquidacionContabilizar)

	return str
}

func Generarasientomanual(w http.ResponseWriter, r *http.Request, cuentasImportes []strCuentaImporte, tokenAutenticacion *structAutenticacion.Security, descripcion string, id string, options string, codigo string) *StrDatosAsientoContableManual {

	str := requestMonoliticoContabilizarLiquidaciones(w, r, cuentasImportes, tokenAutenticacion, descripcion, id, options, codigo)

	var datosAsientoContableManual StrDatosAsientoContableManual

	json.Unmarshal([]byte(str), &datosAsientoContableManual)

	return &datosAsientoContableManual

}

func Checkgeneroasientomanual(w http.ResponseWriter, r *http.Request, cuentasImportes []strCuentaImporte, tokenAutenticacion *structAutenticacion.Security, descripcion string, id string, options string, codigo string) bool {

	//datosAsientoContableManual := Generarasientomanual(w, r, cuentasImportes, tokenAutenticacion, descripcion, id, options, codigo)
	return true
}
