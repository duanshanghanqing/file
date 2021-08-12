;(function(){
    //使用ajax文件上传
    function ajax_upload(obj = {}) {
        if (typeof obj !== 'object') {
            return;
        }

        if(!obj.url) {
            console.warn("url 不可为空");
            return;
        }

        if (!obj.files || obj.files.length === 0) {
            return;
        }
        // 转数组
        obj.files = Array.prototype.slice.call(obj.files);

        obj = Object.assign({
            formName: 'file',
            timeout: 1000,
        }, obj);

        let index = 0;
        let xhr;
        function fileUpload(){
            xhr = new XMLHttpRequest();
            // 1.准备FormData
            let fd = new FormData();
            fd.append(obj.formName, obj.files[index]);
            // 添加数据
            let data = {};
            if (typeof obj.data === "object") {
                data = obj.data;
            } else if (typeof obj.data === "function") {
                data = obj.data();
            }
            Object.keys(data).forEach(function (name) {
                fd.append(name, data[name]);
            });


            // 监听状态，实时响应
            // xhr 和 xhr.upload 都有progress事件，xhr.progress是下载进度，xhr.upload.progress是上传进度
            //这里监听上传进度事件，文件在上次的过程中，会多次触发该事件，返回一个event事件对象
            xhr.upload.onprogress = function(event) {
                obj.uploadedBeing && obj.uploadedBeing(event);
            };

            // 传输开始事件
            xhr.onloadstart = function(event) {
                obj.uploaduStart && obj.uploaduStart(event);
            };
            // xhr.abort();//调用该方法停止ajax上传，停止当前的网络请求

            //每个文件上传成功
            xhr.onload = function(event) {
                if (xhr.responseText) {
                    let res = {};
                    try {
                        res = JSON.parse(xhr.responseText);
                    } catch (e) {}
                    if (res && res.code === 0) {
                        // 上传成功触发
                        obj.uploadSuccess && obj.uploadSuccess(event);
                    }
                }
                // 上传成功失败都触发
                obj.serviceCallback && obj.serviceCallback(xhr.responseText);
                setTimeout(function(){
                    index++
                    if(index < obj.files.length) {
                        fileUpload();
                    }
                    if (index === obj.files.length) {
                        // 上传完成触发
                        obj.finish && obj.finish(index);
                    }
                }, obj.timeout);
            };

            // ajax过程发生错误事件
            xhr.onerror = function(event) {
                obj.uploadError && obj.uploadError(event);
            };

            // ajax被取消，文件上传被取消，说明调用了 xhr.abort();  方法，所触发的事件
            xhr.onabort = function(event) {
                obj.uploadCancelled && obj.uploadCancelled(event);
            };

            // loadend传输结束，不管成功失败都会被触发
            xhr.onloadend = function (event) {
                obj.uploadEnd && obj.uploadEnd(event);
            };

            // 发起ajax请求传送数据
            xhr.open('POST',obj.url , false);

            // 设置头
            let headers = {};
            if (typeof obj.headers === 'function') {
                headers = obj.headers();
            } else if (typeof obj.headers === 'object') {
                headers = obj.headers;
            }
            Object.keys(headers).forEach(function (name) {
                xhr.setRequestHeader(name, headers[name]);
            });

            xhr.send(fd);//发送文件
        }
        fileUpload(index);

        return {
            abort: function() {// ajax被取消，文件上传被取消
                xhr.abort(); // 调用该方法停止ajax上传，停止当前的网络请求
            }
        };
    }

    window.multiFileUpload = ajax_upload
})();