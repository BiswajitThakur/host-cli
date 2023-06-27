package main

/*
// List of blocked ads & web
    fetch("/show_hosts_list").then(res => res.json()).then(data => {
        console.log(data)
    }).catch(err => {
        console.error(err);
    });
*/

/*
// List of unblock web
    fetch("/show_allow_list").then(res => res.json()).then(data => {
        console.log(data)
    }).catch(err => {
        console.error(err);
    });
*/

/*
// List of blocked web, which blocked by you
    fetch("/show_block_list").then(res => res.json()).then(data => {
        console.log(data)
    }).catch(err => {
        console.error(err);
    });
*/

/*
//  Add new www.google.com, www.youtube.com, fooo.com websites to block list
fetch("/add_block_list", {
    method: 'POST',
    body: JSON.stringify(["www.google.com", "www.youtube.com", "fooo.com"]),
    headers: {
        'Content-Type': 'application/json'
    }
}).then(response=>response.json()).then(data=>console.log(data)).catch(error=>console.error(error));
*/

/*
//  Unblock www.google.com, www.youtube.com, fooo.com and these websites will add to allow or unblock list
fetch("/add_allow_list", {
    method: 'POST',
    body: JSON.stringify(["www.google.com", "www.youtube.com", "fooo.com"]),
    headers: {
        'Content-Type': 'application/json'
    }
}).then(response=>response.json()).then(data=>console.log(data)).catch(error=>console.error(error));
*/

/*
// Unblock all ads & web. Note: All blocked webs and ads will not add to allow list.
    fetch("/unblock_ads_and_web").then(res => res.json()).then(data => {
        console.log(data)
    }).catch(err => {
        console.error(err);
    });
*/

/*
// Blocked all ads
    fetch("/block_by_source_list").then(res => res.json()).then(data => {
        console.log(data)
    }).catch(err => {
        console.error(err);
    });
*/

/*
// If you want open www.google.com, www.github.com will be open
// If you want open goo.com, foo.com will be open
fetch("/add_redirect_list", {
    method: 'POST',
    body: JSON.stringify([
        ["www.github.com", "www.google.com"],
        ["foo.com", "goo.com"]
    ]),
    headers: {
        'Content-Type': 'application/json'
    }
}).then(response=>response.json()).then(data=>console.log(data)).catch(error=>console.error(error));
*/

/*
// Remove www.google.com, goo.com from redirect list
fetch("/remove_from_redirec", {
    method: 'POST',
    body: JSON.stringify([
        ["www.github.com", "www.google.com"],
        ["foo.com", "goo.com"]
    ]),
    headers: {
        'Content-Type': 'application/json'
    }
}).then(response=>response.json()).then(data=>console.log(data)).catch(error=>console.error(error));
*/

/*
// remove kkkk.com, 012.in from allow list and these hosts will not add to block list
fetch("/remove_from_allow_list", {
    method: 'POST',
    body: JSON.stringify(["kkkk.com", "012.in"]),
    headers: {
        'Content-Type': 'application/json'
    }
}).then(response=>response.json()).then(data=>console.log(data)).catch(error=>console.error(error));
*/

/*
// remove "hoho.in", "foofoo.com" from block list and these hosts will not add to allow list
fetch("/remove_from_block_list", {
    method: 'POST',
    body: JSON.stringify(["hoho.in", "foofoo.com"]),
    headers: {
        'Content-Type': 'application/json'
    }
}).then(response=>response.json()).then(data=>console.log(data)).catch(error=>console.error(error));
*/

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Respns struct {
	Msg    interface{} `json:"msg"`
	Err    interface{} `json:"error"`
	Errors []error     `json:"errors"`
}

//go:embed static/*
var staticFiles embed.FS

func Server(u string) {

	fsys := fs.FS(staticFiles)
	html, _ := fs.Sub(fsys, "static")
	http.Handle("/", http.FileServer(http.FS(html)))

	// http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/show_hosts_list", hostlist)
	http.HandleFunc("/show_allow_list", allowListS)
	http.HandleFunc("/show_block_list", blockListS)
	http.HandleFunc("/show_sources_list", sourceListS)
	http.HandleFunc("/show_redirect_list", redirectListS)
	http.HandleFunc("/add_allow_list", aal)
	http.HandleFunc("/add_block_list", abl)
	http.HandleFunc("/unblock_ads_and_web", uaw)
	http.HandleFunc("/block_by_source_list", bbsl)
	http.HandleFunc("/add_redirect_list", arl)
	http.HandleFunc("/remove_from_redirec", rfr)
	http.HandleFunc("/remove_from_allow_list", remove_from_allow_list)
	http.HandleFunc("/remove_from_block_list", remove_from_block_list)
	// http.HandleFunc("/add_source_list", asl)
	var v string = u
	var err error
	if regexp.MustCompile(`^\s*\d+\s*$`).MatchString(u) {
		v = fmt.Sprintf("127.0.0.1:%s", v)
	}
	err = http.ListenAndServe(v, nil)
	if err != nil {
		panic(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "host-cli.html")
}

func WriteMsg(m Respns, w http.ResponseWriter, status int) {
	jsonStr, err := json.Marshal(m)
	if err != nil {
		f, _ := json.Marshal(Respns{
			Err: err,
		})
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write(f)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonStr)
}

// blocked hosts name
func hostlist(w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := json.Marshal(ReadLocalHosts(HostPath))
	if err != nil {
		WriteMsg(Respns{Err: err}, w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

// unblock or allow hosts name
func allowListS(w http.ResponseWriter, r *http.Request) {
	q, _, err := GetList(AllowPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err), http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func blockListS(w http.ResponseWriter, r *http.Request) {
	q, _, err := GetList(BlockPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err), http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func redirectListS(w http.ResponseWriter, r *http.Request) {
	q, err := GetRedirectList(RedirectPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err), http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func sourceListS(w http.ResponseWriter, r *http.Request) {
	q, _, err := GetList(SourcePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", err), http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func aal(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteMsg(Respns{Err: fmt.Sprintf("%s: method not allowed", r.Method)}, w, http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteMsg(Respns{Err: err}, w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var reqBodyList []string
	err = json.Unmarshal(body, &reqBodyList)

	if err != nil {
		WriteMsg(Respns{Err: err}, w, http.StatusBadRequest)
		return
	}
	bl := NewSet()
	for _, i := range reqBodyList {
		bl.Add(i)
	}
	errs := OptionAllowFn(bl)
	if len(errs) == 0 {
		fmt.Println(reqBodyList, "added to allow list")
		WriteMsg(Respns{Msg: "success"}, w, http.StatusOK)
		return
	}
	WriteMsg(Respns{Errors: errs}, w, http.StatusInternalServerError)
}

func abl(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteMsg(Respns{Err: fmt.Sprintf("%s: method not allowed", r.Method)}, w, http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteMsg(Respns{Err: err}, w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var reqBodyList []string
	err = json.Unmarshal(body, &reqBodyList)
	if err != nil {
		WriteMsg(Respns{Err: err}, w, http.StatusBadRequest)
		return
	}
	bl := NewSet()
	for _, i := range reqBodyList {
		bl.Add(i)
	}
	errs := OptionBlock(bl, false)
	if len(errs) == 0 {
		fmt.Println(reqBodyList, "added to block list")
		WriteMsg(Respns{Msg: "success"}, w, http.StatusOK)
		return
	}
	WriteMsg(Respns{Errors: errs}, w, http.StatusInternalServerError)
}

func uaw(w http.ResponseWriter, r *http.Request) {
	errs := OptionAllowFn(NewSet())
	if len(errs) == 0 {
		WriteMsg(Respns{Msg: "success"}, w, http.StatusOK)
		return
	}
	WriteMsg(Respns{Errors: errs}, w, http.StatusInternalServerError)
}

func bbsl(w http.ResponseWriter, r *http.Request) {
	errs := OptionBlock(NewSet(), true)
	if len(errs) == 0 {
		WriteMsg(Respns{Msg: "success"}, w, http.StatusOK)
		return
	}
	WriteMsg(Respns{Errors: errs}, w, http.StatusInternalServerError)
}

func arl(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteMsg(Respns{Err: fmt.Sprintf("%s: method not allowed", r.Method)}, w, http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteMsg(Respns{Err: err}, w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var reqBodyList [][2]string
	err = json.Unmarshal(body, &reqBodyList)
	// fmt.Println(reqBodyList)
	if err != nil {
		WriteMsg(Respns{Err: err}, w, http.StatusBadRequest)
		return
	}
	bl := NewSet()
	rdList, err := GetRedirectList(RedirectPath)
	if err != nil {
		WriteMsg(Respns{Err: err}, w, http.StatusInternalServerError)
		return
	}
	for _, i := range rdList {
		bl.AddRedirect(i[0], i[1])
	}
	for _, i := range reqBodyList {
		bl.AddRedirect(i[0], i[1])
	}
	err = WriteRedirectList(RedirectPath, bl)
	if err != nil {
		WriteMsg(Respns{Err: err}, w, http.StatusInternalServerError)
		return
	}
	err = WriteRedirectListHosts(HostPath, bl)
	if err != nil {
		WriteMsg(Respns{Err: err}, w, http.StatusInternalServerError)
		return
	}
	WriteMsg(Respns{Msg: "success"}, w, http.StatusOK)
}

func rfr(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteMsg(Respns{Err: fmt.Sprintf("%s: method not allowed", r.Method)}, w, http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteMsg(Respns{Err: err}, w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var reqBodyList [][2]string
	err = json.Unmarshal(body, &reqBodyList)
	// fmt.Println(reqBodyList)
	if err != nil {
		WriteMsg(Respns{Err: err}, w, http.StatusBadRequest)
		return
	}
	bl := NewSet()
	for _, i := range reqBodyList {
		bl.AddRedirect(i[0], i[1])
	}
	err = RemoveFromRedirect(bl)
	if err != nil {
		WriteMsg(Respns{Err: err}, w, http.StatusInternalServerError)
		return
	}
	WriteMsg(Respns{Msg: "success"}, w, http.StatusOK)
}

func remove_from_allow_list(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteMsg(Respns{Err: fmt.Sprintf("%s: method not allowed", r.Method)}, w, http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteMsg(Respns{Err: err}, w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var reqBodyList []string
	err = json.Unmarshal(body, &reqBodyList)
	if err != nil {
		WriteMsg(Respns{Err: err}, w, http.StatusBadRequest)
		return
	}
	v := NewSet()
	for _, i := range reqBodyList {
		v.Add(i)
	}
	err = RemoveFromAllowList(v)
	if err != nil {
		WriteMsg(Respns{Err: err}, w, http.StatusBadRequest)
		return
	}
	WriteMsg(Respns{Msg: "success"}, w, http.StatusOK)
}

func remove_from_block_list(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		WriteMsg(Respns{Err: fmt.Sprintf("%s: method not allowed", r.Method)}, w, http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteMsg(Respns{Err: err}, w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var reqBodyList []string
	err = json.Unmarshal(body, &reqBodyList)
	if err != nil {
		WriteMsg(Respns{Err: err}, w, http.StatusBadRequest)
		return
	}
	v := NewSet()
	for _, i := range reqBodyList {
		v.Add(i)
	}
	err = RemoveFromBlockList(v)
	if err != nil {
		WriteMsg(Respns{Err: err}, w, http.StatusBadRequest)
		return
	}
	WriteMsg(Respns{Msg: "success"}, w, http.StatusOK)
}
