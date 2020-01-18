(function () {
  'use strict';

  window.App = {
    conn: null,
    cm: null
  };

  App.cm = CodeMirror.fromTextArea(document.getElementById('code'), {
    lineNumbers: true,
    readOnly: 'nocursor',
    mode: 'markdown',
    gitHubSpice: true,
    highlightFormatting: true
  });

  $('#join-btn').click(function (evt) {
    evt.preventDefault();
    $(this).attr({disabled: true});
    var $username = $('#join-form input[name=username]');
    $username.attr({disabled: true});
    App.conn.send('join', { username: $username.val() });
    $('#leave-btn').attr({disabled:false});
  });
  $('#leave-btn').click(function (evt) {
    evt.preventDefault();
    conn.close();
    
  });

  var url = [location.protocol.replace('http', 'ws'), '//', location.host, '/ws'].join('');
  var conn = App.conn = new SocketConnection(url);

  conn.on('open', function () {
    console.log("open");
    $('#conn-status').text('已连接');
    $('#join-btn').attr({ disabled: false});
    $('#leave-btn').attr({disabled: false});
  });

  conn.on('close', function (evt) {
    console.log("close");
    $('#conn-status').text('和服务器失去联系');
    $('#join-btn').attr({ disabled: false});
    var $username = $('#join-form input[name=username]');
    $username.attr({disabled: false});

    $('#leave-btn').attr({disabled: true});
  });

  // 首次接受文档内容
  conn.on('doc', function(data) {
    console.log("doc");
    App.cm.setOption('readOnly', false);
    App.cm.setValue(data.document);
    var serverAdapter = new ot.SocketConnectionAdapter(conn);
    var editorAdapter = new ot.CodeMirrorAdapter(App.cm);
    App.client = new ot.EditorClient(data.revision, data.clients, serverAdapter, editorAdapter);
    var a='';
    $.each(data.clients, function(key,value){
        a = a + ','+value.name;
    });


    $('#users').text(a.substring(1,a.length));
    $('#leave-btn').attr({disable:false});
  });

  // 连接到websocket服务器
  conn.on('registered', function(clientId) {
    App.cm.setOption('readOnly', false);

  });
  // 加入到协作之中
  conn.on('join', function(data) {
    console.log("join");
    var a='';
    $.each(data.clients, function(key,value){
      a = a + ','+value.name;
    });

    $('#users').text(a.substring(1,a.length));


    console.log(data);
  });

  // 有用户退出
  conn.on('quit', function(data) {
    console.log('quit',data);
  });
}());
