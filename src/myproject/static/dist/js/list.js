(function (jq) {
    var CONFIG = {};


    String.prototype.format = function (args) {
        return this.replace(/\{(\w+)\}/g, function (s, i) {
            return args[i];
        });
    };


    function initHearder(){
    $("tbody").empty();
        var header=CONFIG.TABLE.HEADER;
        var tr = document.createElement("tr");
        var td = document.createElement("td");
        $(td).html('<input name="btSelectAll" type="checkbox">').appendTo(tr);
        for (index in header){
            var td = document.createElement("td");
        $(td).html(header[index]).appendTo(tr);
        }
        $("tbody").append(tr);
    };


    function initPage(pageMap){
    $("#page").empty();
    /*
    currpage 1
    firstpage 1
    lastpage 2
    pages (4) [1, 2, 3, 4]
    totalpages 4
    */
    var pageHtml = "";

    pageHtml += '<li name="page" value="'+pageMap.firstpage+'"><a href="#">«</a></li>';
    if (pageMap.pages[0] > pageMap.currpage ) {
        pageHtml += '<li name="page" value="'+pageMap.currpage+'" class="active"><a href="#">'+pageMap.currpage+'</a></li>';
        pageMap.pages.pop();
    };
    $.each(pageMap.pages,function(index,pageNum){
        if (pageNum == pageMap.currpage){
        pageHtml += '<li name="page" value="'+pageNum+'" class="active"><a href="#">'+pageNum+'</a></li>';
        } else {
        pageHtml += '<li name="page" value="'+pageNum+'"><a href="#">'+pageNum+'</a></li>';
        };
    });
    pageHtml+='<li name="page" value="'+pageMap.lastpage+'"><a href="#">»</a></li>';
    pageHtml+='<li><a href="#">共'+pageMap.totalpages+'页</a></li>';
    $("#page").append(pageHtml);
    }


    function getSearchVal (name){
    var username = $("[name='"+name+"']").val();
    return username
    };


    function initBody(page){
        var l = CONFIG.TABLE.BODY;
        var requestDict = {}
        requestDict["page"] = page;
        $.each(CONFIG.SEARCH, function(k,v){
            requestDict[v] = getSearchVal(v);
        });

        $.ajax({ 
            type : "get", 
            url : CONFIG.URL.CONTROL_GET, 
            data: requestDict,
            success : function(result){ 
                console.log(result);
                $.each(result.ObjSlice,function(index,data){
                    var tr = document.createElement("tr");
                    var td = document.createElement("td");
                    $(td).html('<input name="btSelect" value="'+data.Id+'" type="checkbox">').appendTo(tr);
                    $.each(l,function(k,v){
                        var td = document.createElement("td")
                        var type = CONFIG.TABLE.COLUMN[v]["type"]
                        if (type == "text") {
                            $(td).text(data[v]).appendTo(tr);
                        } else if (type=="bool") {

                            console.log(type,data[v],"11111111111111111111111",data,v);
                            var value = data[v];
                            if (value == true) {
                                $(td).html(CONFIG.TABLE.COLUMN[v]["true"]).appendTo(tr);
                            } else if (value == false) {
                                $(td).html(CONFIG.TABLE.COLUMN[v]["false"]).appendTo(tr);
                            }
                        } else if (type == "template") {
                            var d = {};
                            if (CONFIG.TABLE.COLUMN[v]["key"].length > 0 ) {
                                $.each(CONFIG.TABLE.COLUMN[v]["key"],function(i,key){
                                    d[key] = data[key];
                                });
                            }
                            $(td).html(CONFIG.TABLE.COLUMN[v]["template"].format(d)).appendTo(tr);
                        };
                    });
                    $("tbody").append(tr);
                });
                
                initPage(result.PaginatorMap);
                closeShadow();
            }
        }); 
    };


    function initShadow(){
    var index = layer.load(-1, {
        shade: [0.1,'#fff'] //0.1透明度的白色背景
    });
    };


    function closeShadow(){
    layer.closeAll();
    };


    function initTable(page){
    initShadow();
    initHearder();
    initBody(page);
    };


    function bindRefreshButton(){
    $('[name="refresh"]').click(function(){
        initTable(1);
    });
    };


    function getAllCheckedObj(){
    var ids = new Array()
    $("[name='btSelect']").each(function(){
        if ($(this).prop("checked")) {
        var id = $(this).val();
        ids.push(parseInt(id));
        }
        
    });
    return ids;
    };


    function deleteIds(){
    var ids = getAllCheckedObj();
    
    $.ajax({
        url:'{{urlfor "PermissionController.Delete"}}',
        type:"DELETE",
        contentType:"application/json",
        accept : "application/json",
        traditional: true,
        dataType : 'json',
        data:JSON.stringify(ids),
        success:function(result){
        initTable(1);
        },
    });
    };


    function bindDeleteEvent(){
    $("#delete").click(function(){

    //询问框
    layer.confirm(' 您是否要删除所选的项？', {
        btn: ['是','否'], icon: 3, title: '请确认' //按钮
    }, function(){
        console.log('确认');
        layer.closeAll();
        deleteIds();
    }, function(){
        console.log('取消');
    });

    });
    };


    function bindCheckBoxEvent(){
    $(document).on("click","[name='btSelectAll']",function(){
        var allIsCheck = $(this).prop('checked');
        $("[name='btSelect']").each(function(k,v){
        
            $(v).prop('checked',allIsCheck);
        });
        checkDeleteButton();
    });

    $(document).on("click","[name='btSelect']",function(){
        checkDeleteButton(); 
    });
    };


    function checkDeleteButton(){
    var count = 0;
    var checkedCount = 0;
    $("[name='btSelect']").each(function(){
        count +=1;
        var isCheck = $(this).prop('checked')
        if (isCheck == true){
        checkedCount+=1;
        };

    });
    
    if (checkedCount == 0) {
        $("#delete").prop("disabled",true);
    } else {
        $("#delete").prop("disabled",false);
    }
    if (count == checkedCount) {
        $("[name='btSelectAll']").prop("checked",true);
    } else {
        $("[name='btSelectAll']").prop("checked",false);
    };
    };


    function bindCreateEvent(){

    $("#add").bind("click",function(){
        var title = CONFIG.TITLES.ADD_TITLE;
        var url = CONFIG.URL.ADD_CONTROL_GET;
        layer.open({
            type: 2,
            title: title,
            shadeClose: false,
            shade: 0.2,
            maxmin: true,
            shift: 1,
            area: ['1000px', '380px'],
            content: url,
            btn: ['保存', '关闭'],
            yes: function (index, layero) {
                var iframeWin = window[layero.find('iframe')[0]['name']];
                var ret = iframeWin.FormSubmit();
                //提示层
                console.log(ret);
                if (ret.Status) {
                layer.closeAll();
                    layer.msg(ret.Msg);
                    initTable(1);
                } else {
                    layer.msg(ret.Msg, {icon: 5});
                }
                
                
            }
        });
    });
        
    };


    function bindPageEvent(){
    $(document).on("click","body [name='page']",function(){
        $.each($(this),function(k,pageHtml){
            var page = $(pageHtml).prop("value");
            var page = parseInt(page);
            initTable(page);
        });
    });
    };


    function bindSearch() {
    $("#btnSearch").bind("click",function(){
        initTable(1);
    });
    };


    function bindSearchClear(){
    
        $("#btnClearSearch").bind("click",function(){
            $.each(CONFIG.SEARCH, function(k,v){
            $("[name='"+v+"']").val("");
            });
            
            initTable(1);
        });
    };


    function bindEditEvent(){
        $(document).on("click","body [name='edit']",function(){
            var permId = $(this).attr("data");
            permId = parseInt(permId);
            var url = CONFIG.URL.EDIT_CONTROL_GET+permId;
            var title = CONFIG.TITLES.EDIT_TITLE;
            layer.open({
                type: 2,
                title: title,
                shadeClose: false,
                shade: 0.2,
                maxmin: true,
                shift: 1,
                area: ['1000px', '350px'],
                content: url,
                btn: ['保存', '关闭'],
                yes: function (index, layero) {
                    var iframeWin = window[layero.find('iframe')[0]['name']];
                    var ret = iframeWin.FormSubmit();
                    //提示层
                    console.log(ret,111111);
                    if (ret.Status) {
                    layer.closeAll();
                        layer.msg(ret.Msg);
                        initTable(1);
                    } else {
                        layer.msg(ret.Msg, {icon: 5});
                    }
                    
                    
                }
            });

        });
    }


    function bindChangePassword() {
    var title = "修改用户密码";
    $(document).on("click","body [name='changePass']",function(){
        var userId =  $(this).parent().parent().prev().prev().prev().text();
        var url = "/changepass/"+userId
        layer.open({
        type: 2,
        title: title,
        shadeClose: false,
        shade: 0.2,
        maxmin: true,
        shift: 1,
        area: ['400px', '400px'],
        content: url,
        btn: ['保存', '关闭'],
        yes: function (index, layero) {
            var iframeWin = window[layero.find('iframe')[0]['name']];
                var ret = iframeWin.FormSubmit();
            if (ret.Status) {
                layer.closeAll();
                initTable(1);
            } else {
                layer.msg(ret.Msg, {icon: 5});
            }
        }
        });

        });
    
    };


    function bindEvents() {
    bindRefreshButton();
    bindDeleteEvent();
    bindCheckBoxEvent();
    bindCreateEvent();
    bindPageEvent();
    bindSearch();
    bindSearchClear();
    bindEditEvent();
    bindChangePassword();
    };



    


    jq.extend({
        "List":function (config) {
            CONFIG = config;
            initTable(1);
            bindEvents();
        }
    });

})(jQuery);
