'use strict';var m,q=this;function aa(a,b){var c=a.split("."),d=q;c[0]in d||!d.execScript||d.execScript("var "+c[0]);for(var e;c.length&&(e=c.shift());)c.length||void 0===b?d=d[e]?d[e]:d[e]={}:d[e]=b}
function r(a){var b=typeof a;if("object"==b)if(a){if(a instanceof Array)return"array";if(a instanceof Object)return b;var c=Object.prototype.toString.call(a);if("[object Window]"==c)return"object";if("[object Array]"==c||"number"==typeof a.length&&"undefined"!=typeof a.splice&&"undefined"!=typeof a.propertyIsEnumerable&&!a.propertyIsEnumerable("splice"))return"array";if("[object Function]"==c||"undefined"!=typeof a.call&&"undefined"!=typeof a.propertyIsEnumerable&&!a.propertyIsEnumerable("call"))return"function"}else return"null";
else if("function"==b&&"undefined"==typeof a.call)return"object";return b}function ba(a){var b=r(a);return"array"==b||"object"==b&&"number"==typeof a.length}function ca(a){return"string"==typeof a}function da(a){return a[fa]||(a[fa]=++ga)}var fa="closure_uid_"+(1E9*Math.random()>>>0),ga=0;function ha(a,b,c){return a.call.apply(a.bind,arguments)}
function ia(a,b,c){if(!a)throw Error();if(2<arguments.length){var d=Array.prototype.slice.call(arguments,2);return function(){var c=Array.prototype.slice.call(arguments);Array.prototype.unshift.apply(c,d);return a.apply(b,c)}}return function(){return a.apply(b,arguments)}}function u(a,b,c){u=Function.prototype.bind&&-1!=Function.prototype.bind.toString().indexOf("native code")?ha:ia;return u.apply(null,arguments)}var ja=Date.now||function(){return+new Date};
function v(a,b){function c(){}c.prototype=b.prototype;a.V=b.prototype;a.prototype=new c};var ka="closure_listenable_"+(1E6*Math.random()|0);function la(a){try{return!(!a||!a[ka])}catch(b){return!1}}var ma=0;function na(a,b,c,d,e){this.p=a;this.Q=null;this.src=b;this.type=c;this.capture=!!d;this.K=e;this.key=++ma;this.t=this.J=!1}function oa(a){a.t=!0;a.p=null;a.Q=null;a.src=null;a.K=null};var pa="constructor hasOwnProperty isPrototypeOf propertyIsEnumerable toLocaleString toString valueOf".split(" ");function qa(a,b){for(var c,d,e=1;e<arguments.length;e++){d=arguments[e];for(c in d)a[c]=d[c];for(var f=0;f<pa.length;f++)c=pa[f],Object.prototype.hasOwnProperty.call(d,c)&&(a[c]=d[c])}};var x=Array.prototype,sa=x.indexOf?function(a,b,c){return x.indexOf.call(a,b,c)}:function(a,b,c){c=null==c?0:0>c?Math.max(0,a.length+c):c;if(ca(a))return ca(b)&&1==b.length?a.indexOf(b,c):-1;for(;c<a.length;c++)if(c in a&&a[c]===b)return c;return-1},ta=x.forEach?function(a,b,c){x.forEach.call(a,b,c)}:function(a,b,c){for(var d=a.length,e=ca(a)?a.split(""):a,f=0;f<d;f++)f in e&&b.call(c,e[f],f,a)};function ua(a){var b=a.length;if(0<b){for(var c=Array(b),d=0;d<b;d++)c[d]=a[d];return c}return[]};function va(a){this.src=a;this.b={};this.c=0}va.prototype.add=function(a,b,c,d,e){var f=this.b[a];f||(f=this.b[a]=[],this.c++);var g=wa(f,b,d,e);-1<g?(a=f[g],c||(a.J=!1)):(a=new na(b,this.src,a,!!d,e),a.J=c,f.push(a));return a};function xa(a,b){var c=b.type;if(c in a.b){var d=a.b[c],e=sa(d,b),f;(f=0<=e)&&x.splice.call(d,e,1);f&&(oa(b),0==a.b[c].length&&(delete a.b[c],a.c--))}}function wa(a,b,c,d){for(var e=0;e<a.length;++e){var f=a[e];if(!f.t&&f.p==b&&f.capture==!!c&&f.K==d)return e}return-1};var y,ya,za,Aa;function Ba(){return q.navigator?q.navigator.userAgent:null}Aa=za=ya=y=!1;var Ca;if(Ca=Ba()){var Da=q.navigator;y=0==Ca.lastIndexOf("Opera",0);ya=!y&&(-1!=Ca.indexOf("MSIE")||-1!=Ca.indexOf("Trident"));za=!y&&-1!=Ca.indexOf("WebKit");Aa=!y&&!za&&!ya&&"Gecko"==Da.product}var Ea=y,z=ya,A=Aa,Fa=za;function Ga(){var a=q.document;return a?a.documentMode:void 0}var Ha;
a:{var Ia="",Ja;if(Ea&&q.opera)var Ka=q.opera.version,Ia="function"==typeof Ka?Ka():Ka;else if(A?Ja=/rv\:([^\);]+)(\)|;)/:z?Ja=/\b(?:MSIE|rv)[: ]([^\);]+)(\)|;)/:Fa&&(Ja=/WebKit\/(\S+)/),Ja)var La=Ja.exec(Ba()),Ia=La?La[1]:"";if(z){var Ma=Ga();if(Ma>parseFloat(Ia)){Ha=String(Ma);break a}}Ha=Ia}var Na={};
function B(a){var b;if(!(b=Na[a])){b=0;for(var c=String(Ha).replace(/^[\s\xa0]+|[\s\xa0]+$/g,"").split("."),d=String(a).replace(/^[\s\xa0]+|[\s\xa0]+$/g,"").split("."),e=Math.max(c.length,d.length),f=0;0==b&&f<e;f++){var g=c[f]||"",h=d[f]||"",k=/(\d*)(\D*)/g,l=/(\d*)(\D*)/g;do{var p=k.exec(g)||["","",""],n=l.exec(h)||["","",""];if(0==p[0].length&&0==n[0].length)break;b=((0==p[1].length?0:parseInt(p[1],10))<(0==n[1].length?0:parseInt(n[1],10))?-1:(0==p[1].length?0:parseInt(p[1],10))>(0==n[1].length?
0:parseInt(n[1],10))?1:0)||((0==p[2].length)<(0==n[2].length)?-1:(0==p[2].length)>(0==n[2].length)?1:0)||(p[2]<n[2]?-1:p[2]>n[2]?1:0)}while(0==b)}b=Na[a]=0<=b}return b}var Oa=q.document,Pa=Oa&&z?Ga()||("CSS1Compat"==Oa.compatMode?parseInt(Ha,10):5):void 0;var Qa=!z||z&&9<=Pa,Ra=z&&!B("9");!Fa||B("528");A&&B("1.9b")||z&&B("8")||Ea&&B("9.5")||Fa&&B("528");A&&!B("8")||z&&B("9");function C(){0!=Sa&&(Ta[da(this)]=this)}var Sa=0,Ta={};C.prototype.j=!1;C.prototype.$=function(){if(!this.j&&(this.j=!0,this.n(),0!=Sa)){var a=da(this);delete Ta[a]}};C.prototype.n=function(){if(this.r)for(;this.r.length;)this.r.shift()()};function D(a,b){this.type=a;this.b=this.u=b}D.prototype.$=function(){};D.prototype.c=!1;D.prototype.e=!0;D.prototype.preventDefault=function(){this.e=!1};function Ua(a){Ua[" "](a);return a}Ua[" "]=function(){};function Va(a,b){if(a){var c=this.type=a.type;D.call(this,c);this.u=a.target||a.srcElement;this.b=b;var d=a.relatedTarget;if(d&&A)try{Ua(d.nodeName)}catch(e){}this.offsetX=Fa||void 0!==a.offsetX?a.offsetX:a.layerX;this.offsetY=Fa||void 0!==a.offsetY?a.offsetY:a.layerY;this.keyCode=a.keyCode||0;this.charCode=a.charCode||("keypress"==c?a.keyCode:0);this.ctrlKey=a.ctrlKey;this.altKey=a.altKey;this.aa=a;a.defaultPrevented&&this.preventDefault();delete this.c}}v(Va,D);m=Va.prototype;m.u=null;
m.offsetX=0;m.offsetY=0;m.keyCode=0;m.charCode=0;m.ctrlKey=!1;m.altKey=!1;m.aa=null;m.preventDefault=function(){Va.V.preventDefault.call(this);var a=this.aa;if(a.preventDefault)a.preventDefault();else if(a.returnValue=!1,Ra)try{if(a.ctrlKey||112<=a.keyCode&&123>=a.keyCode)a.keyCode=-1}catch(b){}};var Wa={},Xa={},E={};function Ya(a,b,c,d,e){if("array"==r(b))for(var f=0;f<b.length;f++)Ya(a,b[f],c,d,e);else if(c=Za(c),la(a))a.m.add(b,c,!1,d,e);else{f=c;if(!b)throw Error("Invalid event type");c=!!d;var g=da(a),h=Xa[g];h||(Xa[g]=h=new va(a));d=h.add(b,f,!1,d,e);d.Q||(e=bb(),d.Q=e,e.src=a,e.p=d,a.addEventListener?a.addEventListener(b,e,c):a.attachEvent(b in E?E[b]:E[b]="on"+b,e),Wa[d.key]=d)}}
function bb(){var a=cb,b=Qa?function(c){return a.call(b.src,b.p,c)}:function(c){c=a.call(b.src,b.p,c);if(!c)return c};return b}function db(a,b,c,d,e){if("array"==r(b))for(var f=0;f<b.length;f++)db(a,b[f],c,d,e);else(c=Za(c),la(a))?(a=a.m,b in a.b&&(f=a.b[b],c=wa(f,c,d,e),-1<c&&(oa(f[c]),x.splice.call(f,c,1),0==f.length&&(delete a.b[b],a.c--)))):a&&(d=!!d,a=eb(a))&&(b=a.b[b],a=-1,b&&(a=wa(b,c,d,e)),(c=-1<a?b[a]:null)&&fb(c))}
function fb(a){if("number"!=typeof a&&a&&!a.t){var b=a.src;if(la(b))xa(b.m,a);else{var c=a.type,d=a.Q;b.removeEventListener?b.removeEventListener(c,d,a.capture):b.detachEvent&&b.detachEvent(c in E?E[c]:E[c]="on"+c,d);(c=eb(b))?(xa(c,a),0==c.c&&(c.src=null,delete Xa[da(b)])):oa(a);delete Wa[a.key]}}}function gb(a,b,c,d){var e=1;if(a=eb(a))if(b=a.b[b])for(b=ua(b),a=0;a<b.length;a++){var f=b[a];f&&f.capture==c&&!f.t&&(e&=!1!==hb(f,d))}return Boolean(e)}
function hb(a,b){var c=a.p,d=a.K||a.src;a.J&&fb(a);return c.call(d,b)}
function cb(a,b){if(a.t)return!0;if(!Qa){var c;if(!(c=b))a:{c=["window","event"];for(var d=q,e;e=c.shift();)if(null!=d[e])d=d[e];else{c=null;break a}c=d}e=c;c=new Va(e,this);d=!0;if(!(0>e.keyCode||void 0!=e.returnValue)){a:{var f=!1;if(0==e.keyCode)try{e.keyCode=-1;break a}catch(g){f=!0}if(f||void 0==e.returnValue)e.returnValue=!0}e=[];for(f=c.b;f;f=f.parentNode)e.push(f);for(var f=a.type,h=e.length-1;!c.c&&0<=h;h--)c.b=e[h],d&=gb(e[h],f,!0,c);for(h=0;!c.c&&h<e.length;h++)c.b=e[h],d&=gb(e[h],f,!1,
c)}return d}return hb(a,new Va(b,this))}function eb(a){return a[fa]?Xa[da(a)]||null:null}var ib="__closure_events_fn_"+(1E9*Math.random()>>>0);function Za(a){return"function"==r(a)?a:a[ib]||(a[ib]=function(b){return a.handleEvent(b)})};var jb="StopIteration"in q?q.StopIteration:Error("StopIteration");function kb(){}kb.prototype.b=function(){throw jb;};kb.prototype.W=function(){return this};function lb(a){if(a instanceof kb)return a;if("function"==typeof a.W)return a.W(!1);if(ba(a)){var b=0,c=new kb;c.b=function(){for(;;){if(b>=a.length)throw jb;if(b in a)return a[b++];b++}};return c}throw Error("Not implemented");}
function mb(a,b){if(ba(a))try{ta(a,b,void 0)}catch(c){if(c!==jb)throw c;}else{a=lb(a);try{for(;;)b.call(void 0,a.b(),void 0,a)}catch(d){if(d!==jb)throw d;}}};function F(a,b){this.o={};this.g=[];this.c=this.b=0;var c=arguments.length;if(1<c){if(c%2)throw Error("Uneven number of arguments");for(var d=0;d<c;d+=2)nb(this,arguments[d],arguments[d+1])}else if(a){var e;if(a instanceof F)for(d=ob(a),pb(a),e=[],c=0;c<a.g.length;c++)e.push(a.o[a.g[c]]);else{var c=[],f=0;for(d in a)c[f++]=d;d=c;c=[];f=0;for(e in a)c[f++]=a[e];e=c}for(c=0;c<d.length;c++)nb(this,d[c],e[c])}}function ob(a){pb(a);return a.g.concat()}
function pb(a){if(a.b!=a.g.length){for(var b=0,c=0;b<a.g.length;){var d=a.g[b];Object.prototype.hasOwnProperty.call(a.o,d)&&(a.g[c++]=d);b++}a.g.length=c}if(a.b!=a.g.length){for(var e={},c=b=0;b<a.g.length;)d=a.g[b],Object.prototype.hasOwnProperty.call(e,d)||(a.g[c++]=d,e[d]=1),b++;a.g.length=c}}function nb(a,b,c){Object.prototype.hasOwnProperty.call(a.o,b)||(a.b++,a.g.push(b),a.c++);a.o[b]=c}F.prototype.e=function(){return new F(this)};
F.prototype.W=function(a){pb(this);var b=0,c=this.g,d=this.o,e=this.c,f=this,g=new kb;g.b=function(){for(;;){if(e!=f.c)throw Error("The map has changed since the iterator was created");if(b>=c.length)throw jb;var g=c[b++];return a?g:d[g]}};return g};function G(a,b){C.call(this);this.f=b;this.c=[];if(a>this.f)throw Error("[goog.structs.SimplePool] Initial cannot be greater than max");for(var c=0;c<a;c++)this.c.push(this.b())}v(G,C);G.prototype.b=function(){return{}};G.prototype.e=function(a){var b=typeof a;if("object"==b&&null!=a||"function"==b)if("function"==r(a.$))a.$();else for(var c in a)delete a[c]};G.prototype.n=function(){G.V.n.call(this);for(var a=this.c;a.length;)this.e(a.pop());delete this.c};function qb(){this.c=[];this.f=new F;this.b=new F;this.r=1;this.j=new G(0,4E3);this.j.b=function(){return new rb};this.v=new G(0,50);this.v.b=function(){return new sb};var a=this;this.e=new G(0,2E3);this.e.b=function(){return String(a.r++)};this.e.e=function(){}}function sb(){this.time=this.count=0}sb.prototype.toString=function(){var a=[];a.push(this.type," ",this.count," (",Math.round(10*this.time)/10," ms)");return a.join("")};function rb(){}
function tb(a,b,c){var d=[];-1==b?d.push("    "):d.push(ub(a.c-b));d.push(" ",vb(a.c-0));0==a.b?d.push(" Start        "):1==a.b?(d.push(" Done "),d.push(ub(a.j-a.startTime)," ms ")):d.push(" Comment      ");d.push(c,a);0<a.f&&d.push("[VarAlloc ",a.f,"] ");return d.join("")}rb.prototype.toString=function(){return null==this.type?this.e:"["+this.type+"] "+this.e};
qb.prototype.toString=function(){for(var a=[],b=-1,c=[],d=0;d<this.c.length;d++){var e=this.c[d];1==e.b&&c.pop();a.push(" ",tb(e,b,c.join("")));b=e.c;a.push("\n");0==e.b&&c.push("|  ")}if(0!=this.f.b){var f=ja();a.push(" Unstopped timers:\n");mb(this.f,function(b){a.push("  ",b," (",f-b.startTime," ms, started at ",vb(b.startTime),")\n")})}b=ob(this.b);for(d=0;d<b.length;d++)c=Object.prototype.hasOwnProperty.call(this.b.o,b[d])?this.b.o[b[d]]:void 0,1<c.count&&a.push(" TOTAL ",c,"\n");a.push("Total tracers created ",
0,"\n","Total comments created ",0,"\n","Overhead start: ",0," ms\n","Overhead end: ",0," ms\n","Overhead comment: ",0," ms\n");return a.join("")};function ub(a){a=Math.round(a);var b="";1E3>a&&(b=" ");100>a&&(b="  ");10>a&&(b="   ");return b+a}function vb(a){a=Math.round(a);return String(100+a/1E3%60).substring(1,3)+"."+String(1E3+a%1E3).substring(1,4)}new qb;function H(){C.call(this);this.m=new va(this);this.v=this}v(H,C);H.prototype[ka]=!0;H.prototype.f=null;H.prototype.removeEventListener=function(a,b,c,d){db(this,a,b,c,d)};
function wb(a,b){var c,d=a.f;if(d)for(c=[];d;d=d.f)c.push(d);var d=a.v,e=b,f=e.type||e;if(ca(e))e=new D(e,d);else if(e instanceof D)e.u=e.u||d;else{var g=e,e=new D(f,d);qa(e,g)}var g=!0,h;if(c)for(var k=c.length-1;!e.c&&0<=k;k--)h=e.b=c[k],g=xb(h,f,!0,e)&&g;e.c||(h=e.b=d,g=xb(h,f,!0,e)&&g,e.c||(g=xb(h,f,!1,e)&&g));if(c)for(k=0;!e.c&&k<c.length;k++)h=e.b=c[k],g=xb(h,f,!1,e)&&g}
H.prototype.n=function(){H.V.n.call(this);if(this.m){var a=this.m,b=0,c;for(c in a.b){for(var d=a.b[c],e=0;e<d.length;e++)++b,oa(d[e]);delete a.b[c];a.c--}}this.f=null};function xb(a,b,c,d){b=a.m.b[b];if(!b)return!0;b=ua(b);for(var e=!0,f=0;f<b.length;++f){var g=b[f];if(g&&!g.t&&g.capture==c){var h=g.p,k=g.K||g.src;g.J&&xa(a.m,g);e=!1!==h.call(k,d)&&e}}return e&&!1!=d.e};function yb(a,b){H.call(this);this.ca=void 0!==a?a:!0;this.c=b||zb;this.e=this.c(this.F)}v(yb,H);m=yb.prototype;m.k=null;m.H=null;m.B=void 0;m.Z=!1;m.F=0;function zb(a){return Math.min(1E3*Math.pow(2,a),6E4)}m.ba=function(a,b){null!=this.b&&q.clearTimeout(this.b);this.b=null;this.H=a;this.k=(this.B=b)?new WebSocket(this.H,this.B):new WebSocket(this.H);this.k.onopen=u(this.ga,this);this.k.onclose=u(this.da,this);this.k.onmessage=u(this.fa,this);this.k.onerror=u(this.ea,this)};
function Ab(a){null!=a.b&&q.clearTimeout(a.b);a.b=null;a.k&&(a.Z=!0,a.k.close(),a.k=null)}m.ga=function(){wb(this,"d");this.F=0;this.e=this.c(this.F)};
m.da=function(){wb(this,"a");this.k=null;if(this.Z)this.H=null,this.B=void 0;else if(this.ca){var a=u(this.ba,this,this.H,this.B),b=this.e;if("function"==r(a))this&&(a=u(a,this));else if(a&&"function"==typeof a.handleEvent)a=u(a.handleEvent,a);else throw Error("Invalid listener argument");this.b=2147483647<b?-1:q.setTimeout(a,b||0);this.F++;this.e=this.c(this.F)}this.Z=!1};m.fa=function(a){wb(this,new Bb(a.data))};m.ea=function(a){wb(this,new Cb(a.data))};m.n=function(){yb.V.n.call(this);Ab(this)};
function Bb(a){D.call(this,"c");this.message=a}v(Bb,D);function Cb(a){D.call(this,"b");this.data=a}v(Cb,D);var Db=new yb,Eb=!1;function I(a){Eb&&(a=JSON.stringify(a),Db.k.send(a))}Db.ba("wss://"+location.host+"/ws");Ya(Db,"d",Fb);Ya(Db,"a",Gb);Ya(Db,"c",Hb);var Ib=[],Jb=[];function Fb(){Eb=!0;Ib.forEach(function(a){a()})}function Gb(){Eb=!1;Jb.forEach(function(a){a()})}var J={};function Hb(a){a=JSON.parse(a.message);var b,c;for(c in a)if(b=J[c])"Kick"==c&&Ab(Db),b(a[c])};var K=document.createElement("canvas");aa("screenshot",function(){window.open(K.toDataURL())});var L=K.getContext("2d"),M=window.innerWidth,N=window.innerHeight;K.width=M;K.height=N;window.addEventListener("resize",function(){K.width=M=window.innerWidth;K.height=N=window.innerHeight;O()},!1);var Kb=Infinity,Lb=0;function O(a){a=!isNaN(a)&&0<a?a:0;var b=a+Date.now();Kb>b&&(Lb&&clearTimeout(Lb),Lb=setTimeout(function(){window.requestAnimationFrame(Mb);Lb=0;Kb=Infinity},a),Kb=b)}var Nb=!1;
Ib.push(function(){Nb=!0;O()});Jb.push(function(){Nb=!1;O()});
function Mb(){L.clearRect(0,0,M,N);if(Nb){for(var a=M,b=N,c=Ob(),d=Pb(),c=a/2/32-c,d=b/2/32-d,e=0;256>e;e+=16)for(var f=0;256>f;f+=16)Qb.d(c+e+8-0.5,d+f+16-1);for(var g=0;2>g;g++)for(Rb=!g,f=Math.max(0,Math.floor(-b/2-d));f<Math.min(256,Math.floor(b/2-d));f++){for(e=Math.max(0,Math.floor((-a/2-c)/2));e<Math.min(128,Math.floor((a/2-c)/2));e++){var h=Q[e<<1|f<<8];if(h)for(var k in h)h[k].d(c,d)}for(e=Math.max(0,Math.floor((-a/2-c)/2));e<Math.min(128,Math.floor((a/2-c)/2));e++)if(h=Q[e<<1|f<<8|1])for(k in h)h[k].d(c,
d)}Sb()}else Tb.d(M/2/32,N/2/32)}K.onmousemove=function(a){var b=a.offsetX||a.layerX;a=a.offsetY||a.layerY;var c=M,d=N;Vb(b,a,c,d)&&(a=b=-Infinity);Wb(b,a,c,d)};K.onmouseout=function(){var a=M,b=N;Vb(-Infinity,-Infinity,a,b);Wb(-Infinity,-Infinity,a,b)};var Xb=!0;window.addEventListener("blur",function(){Xb=!1},!1);window.addEventListener("focus",function(){setTimeout(function(){Xb=!0},100)},!1);
K.onclick=function(a){if(Xb){var b=a.offsetX||a.layerX,c=a.offsetY||a.layerY,d=M;a=N;var e;a:if(Vb(b,c,d,a),Yb)Zb("inv"),e=!0;else{for(e=R.length-1;0<=e;e--)if(R[e].e(b,c,d,a)){e=!0;break a}e=!1}if(!(e||(d=d/2/32-Ob(),a=a/2/32-Pb(),b=Math.floor(b/32-d),c=Math.ceil(c/32-a),0>b||256<=b||0>c||256<=c))){a=[];for(var f in Q[b|c<<8])"_"!=f.charAt(0)&&(d=Q[b|c<<8][f],"_fl"!=d.b[0].c&&d.l.length&&a.unshift(d));a.length?I({Interact:{ID:a[0].id,X:a[0].x,Y:a[0].y,Action:a[0].l[0]}}):I({Walk:{X:b,Y:c}})}}else Xb=
!0};K.oncontextmenu=function(a){a.preventDefault();Xb=!0;var b=a.offsetX||a.layerX,c=a.offsetY||a.layerY,d=M;a=N;var e;a:{Vb(b,c,d,a);for(e=R.length-1;0<=e;e--)if(R[e].j(b,c,d,a)){e=!0;break a}e=!1}if(!(e||(d=d/2/32-Ob(),a=a/2/32-Pb(),b=Math.floor(b/32-d),c=Math.ceil(c/32-a),0>b||256<=b||0>c||256<=c))){a=[];for(var f in Q[b|c<<8])"_"!=f.charAt(0)&&a.unshift(Q[b|c<<8][f]);1==a.length?Zb("menu",a[0]):a.length&&Zb("menu2",a)}};function $b(a,b,c,d,e,f,g,h){function k(){function c(a,b){return e?a:128<=a?255-(255-a)*(255-b)/127:a*b/127}if(n.width&&n.height){l.width=n.width;l.height=n.height;var d=l.getContext("2d");a in S||(S[a]={});b in S[a]||(S[a][b]=[]);if(p in S[a][b])l.width*=p,l.height*=p,d.drawImage(S[a][b][p],0,0);else{S[a][b][p]=l;var e="no"==b;e&&(b="#000");d.fillStyle=b;d.fillRect(0,0,1,1);var f=d.getImageData(0,0,1,1),g=f.data[0],h=f.data[1],k=f.data[2],Ac=f.data[3];d.clearRect(0,0,1,1);d.drawImage(n,0,0);f=d.getImageData(0,
0,n.width,n.height);l.width*=p;l.height*=p;for(var ra=d.getImageData(0,0,l.width,l.height),Ub=0,P=0,ea=0,$a=0;$a<l.height;$a++){for(var ab=0;ab<l.width;ab++)ra.data[ea+0]=c(f.data[P+0],g),ra.data[ea+1]=c(f.data[P+1],h),ra.data[ea+2]=c(f.data[P+2],k),ra.data[ea+3]=f.data[P+3]*Ac/255,ab%p==p-1&&(P+=4),ea+=4;$a%p==p-1?Ub=P:P=Ub}d.putImageData(ra,0,0);O()}}}var l=document.createElement("canvas");this.j=l;this.c=c;this.r=d;this.v=e;this.f=f;this.e=g;var p=Math.floor(h)||1;this.b=p;var n;a+=".png";(n=ac[a])?
k():(n=new Image,n.onload=function(){ac[a]=n;k()},n.src=a)}var ac={},S={};
$b.prototype.d=function(a,b){function c(){L.drawImage(f.j,Math.floor(d*f.f*f.b),Math.floor(e*f.e*f.b),Math.floor(f.f*f.b),Math.floor(f.e*f.b),Math.floor(32*a-(f.f*f.b-32)/2),Math.floor(32*b-f.e*f.b),Math.floor(f.f*f.b),Math.floor(f.e*f.b))}var d=this.r,e=this.v,f=this;if(Rb)"_fl"==this.c&&(b+=0.5,c());else{switch(this.c){case "ccr":d+=[0,6,3,9][Math.floor(Date.now()/1500)%4];O(1500-Date.now()%1500);break;case "wa":d+=[0,1,0,2][Math.floor(Date.now()/150)%4];O(150-Date.now()%150);break;case "l2":d+=
Math.floor(Date.now()/150)%2;O(150-Date.now()%150);break;case "l3":d+=Math.floor(Date.now()/150)%3;O(150-Date.now()%150);break;case "_ac":if(2>e)break;var g=Date.now()/1E4;switch(e){case 2:case 3:case 4:case 5:b+=Math.sin(5*g+7*Math.cos(3*g)+e)/8}O(100);break;case "wa_ac":d+=[0,1,0,2][Math.floor(Date.now()/150)%4];O(50);if(2>e)break;g=Date.now()/1E4;switch(e){case 2:case 3:case 4:case 5:b+=Math.sin(5*g+7*Math.cos(3*g)+e)/8}break;case "_fl":return}c()}};var Rb=!1;
function bc(a){return new $b(a.S,a.C,a.E.a||"",a.E.x||0,a.E.y||0,a.E.w||32,a.E.h||32,a.E.s||1)};function cc(a,b,c,d){this.id=c;this.L=this.x=a;this.M=this.y=b;this.c=0;this.name=d.N;var e=[];this.b=e;(d.S||[]).forEach(function(a){e.push(bc(a))});this.l=d.A||[]}cc.prototype.e=function(){var a=new cc(this.x,this.y,this.id,{N:this.name,A:this.l,S:null});a.b=this.b;return a};function dc(a,b){a.name=b.N;a.l=b.A||[];var c=a.b;c.length=0;(b.S||[]).forEach(function(a){c.push(bc(a))})}function ec(a,b,c){a.L=fc(a.x,a.L,a.c);a.M=fc(a.y,a.M,a.c);a.x=b;a.y=c;a.c=Date.now()}
cc.prototype.d=function(a,b){var c=fc(this.x,this.L,this.c),d=fc(this.y,this.M,this.c);this.b.forEach(function(e){e.d(c+a,d+b)})};function fc(a,b,c){c=Date.now()-c;return 400>c?(O(),(c*a+(400-c)*b)/400):a};function T(a,b,c,d){this.b=a;this.j=b;this.e=c?1:0.5;this.c=32*this.e+"px "+(c?'"Jolly Lodger"':'"Open Sans Condensed"');this.f=!d}T.prototype.width=function(){L.font=this.c;return L.measureText(this.b).width/32};T.prototype.height=function(){return this.e};
T.prototype.d=function(a,b){L.font=this.c;var c=L.measureText(this.b).width;a=Math.floor(32*a-(this.f?c/2:0));b=Math.floor(32*b);L.fillStyle="rgba(0,0,0,.2)";for(c=-1;1>=c;c++)for(var d=-1;2>=d;d++)L.fillText(this.b,a+c,b+d);L.fillStyle=this.j;L.fillText(this.b,a,b)};var Tb=new T("connection lost","#aaa",!0);function gc(a,b,c,d,e){this.c=a;this.f=b;this.b=c;this.e=d;this.j=e}
var R=[],hc=new T("Character","#ccc",!0),ic=new T("change name","#aaa",!1),jc=new T("change name","#fff",!1),kc=new T("change skin color","#aaa",!1),lc=new T("change skin color","#fff",!1),mc=new T("change shirt color","#aaa",!1),nc=new T("change shirt color","#fff",!1),oc=new T("change pants color","#aaa",!1),pc=new T("change pants color","#fff",!1),qc=new T("Confirm","#aaa",!0),rc=new T("Confirm","#fff",!0),sc=new T("[push enter to chat]","#ccc",!1,!0),U={};
function Zb(a,b){var c=U[a](b||{});c.b(tc[0],tc[1],tc[2],tc[3]);var d=!1;R.forEach(function(b,f){if(b.c==a)return d=!0,R[f]=c,!1});d||R.push(c);O()}function V(a){R.forEach(function(b,c){if(b.c==a)return R.splice(c,1),O(),!1})}J.HUD=function(a){R.length=0;O();a.N.length&&Zb(a.N,a.D)};var Yb=!1,uc=new $b("ui_icons","no","",0,0,32,32),vc=new $b("ui_icons","#bbb","",0,0,32,32);
function Sb(){var a=M,b=N;null==W?sc.d(0.1,b/32-0.1):(new T(W+"_","#fff",!1,!0)).d(0.1,b/32-0.1);wc.forEach(function(a,d){a.d(0.1,b/32-d/2-0.6)});Yb?vc.d(a/32-1,b/32):uc.d(a/32-1,b/32);R.forEach(function(c){c.f(a,b)})}var tc=[-Infinity,-Infinity,1,1];function Vb(a,b,c,d){tc=[a,b,c,d];if(a>=c-32&&a<=c&&b>=d-32&&b<=d)return Yb||(Yb=!0,O()),!0;Yb&&(Yb=!1,O());for(var e=R.length-1;0<=e;e--)if(R[e].b(a,b,c,d))return!0;return!1}
window.addEventListener("keydown",function(a){if(K.parentNode&&!a.ctrlKey&&!a.altKey)switch(20>a.keyCode&&a.preventDefault(),a.keyCode){case 8:null!==W&&0<W.length&&(W=W.substring(0,W.length-1),O());break;case 13:null===W?W="":(I({Chat:W}),W=null);O();break;case 27:null!==W&&(W=null,O())}},!1);window.addEventListener("keypress",function(a){null!==W&&(W+=String.fromCharCode(a.charCode),O())},!1);Ib.push(function(){R.length=0;W=null});var W=null,wc=[];
J.Msg=function(a){a.forEach(function(a){wc.unshift(new T(a.T,a.C,!1,!0));window.console.log(a.T);window.setTimeout(function(){wc.pop()},6E4);O()})};var X=[];J.Inventory=function(a){var b={};X.forEach(function(a){b[a.id]=a});X.length=0;a.forEach(function(a,d){var e=b[a.I];e?(dc(e,a.O),ec(e,d%8,Math.floor(d/8))):e=new cc(d%8,Math.floor(d/8),a.I,a.O);X.push(e)});O()};var Y=window.sessionStorage.rnoadm_username||"",Z=window.sessionStorage.rnoadm_password||"",$=document.querySelector("form"),xc=$.querySelector("#username"),yc=$.querySelector("#password"),zc=$.querySelector("#password2");Ib.push(function(){Y.length&&2<Z.length&&(xc.value=Y,yc.value=Z,zc.value=Z,$.onsubmit())});aa("admin",function(a){I({Admin:[].slice.call(arguments).map(String)})});yc.onchange=function(){zc.value=yc.value};
$.onsubmit=function(){Y=xc.value;Z=yc.value;if(Y.length)if(2>=Z.length)Y=Z="",yc.focus();else{window.sessionStorage.rnoadm_username=Y;window.sessionStorage.rnoadm_password=Z;I({Auth:{U:Y,P:Z}});var a=$.parentNode;a&&(a.removeChild($),a.style.overflow="hidden",a.style.fontSize="0",a.appendChild(K))}else Y=Z="",xc.focus()};J.Kick=function(a){window.sessionStorage.rnoadm_username=Y="";window.sessionStorage.rnoadm_password=Z="";alert(a)};var Q=Array(65536);J.Update=function(a){a.forEach(function(a){var c=a.I,d=a.X,e=a.Y,f=a.Fx,g=a.Fy,h=a.R,k=a.O;a=d|e<<8;Q[a]||(Q[a]={});h?(delete Q[a][c],O()):void 0!==k?c in Q[a]?dc(Q[a][c],k):Q[a][c]=new cc(d,e,c,k):(f|=g<<8,Q[a][c]=Q[f][c],delete Q[f][c],ec(Q[a][c],d,e))});O()};var Bc=0,Cc=127,Dc=127;J.PlayerX=function(a){Dc!=a&&(Cc=Ob(),Dc=a,Bc=Date.now(),O())};function Ob(){return fc(Dc,Cc,Bc)}var Ec=0,Fc=127,Gc=127;J.PlayerY=function(a){Gc!=a&&(Fc=Pb(),Gc=a,Ec=Date.now(),O())};
function Pb(){return fc(Gc,Fc,Ec)}Jb.push(function(){for(var a=0;a<Q.length;a++)delete Q[a]});var Qb=new $b("grass","#268f1e","",0,0,512,512),Hc=-1,Ic=-1;
function Wb(a,b,c,d){c=c/2/32-Ob();d=d/2/32-Pb();a=Math.floor(a/32-c);b=Math.ceil(b/32-d);if(a!=Hc||b!=Ic)if(K.style.cursor="",!(0>a||255<a||0>b||255<b)){d="";for(var e in Q[a|b<<8])"_"!=e.charAt(0)&&(c=Q[a|b<<8][e],"_fl"!=c.b[0].c&&c.l.length&&(d=c.l[0]));d&&(K.style.cursor="url(cursor_"+d+".png),auto");-1==Hc&&(Hc=a,Ic=b);I({Mouse:{Fx:Hc,Fy:Ic,X:a,Y:b}});Hc=a;Ic=b}};new function(){ja()};!A&&!z||z&&z&&9<=Pa||A&&B("1.9.1");z&&B("9");var Jc;
U.cc=function(a){var b=[];a.S.forEach(function(a){b.push(bc(a))});var c=new T(a.N,"#aaa",!0),d=new T(a.N,"#fff",!0),e=new T("change gender ("+a.G+")","#aaa",!1),f=new T("change gender ("+a.G+")","#fff",!1),g=!1,h=!1,k=!1,l=!1,p=!1,n=!1,t=!0;return new gc("cc",function(a,s){L.fillStyle="rgba(0,0,0,.7)";L.fillRect(0,0,a,s);a/=64;s/=64;b.forEach(function(b){b.d(a-3,s+1)});hc.d(a+2,s-4);g?(d.d(a,s+2),jc.d(a,s+2.5)):(c.d(a,s+2),ic.d(a,s+2.5));h?f.d(a+2,s-3):e.d(a+2,s-3);k?lc.d(a+2,s-2):kc.d(a+2,s-2);l?
nc.d(a+2,s-1):mc.d(a+2,s-1);p?pc.d(a+2,s):oc.d(a+2,s);n?rc.d(a,s+4):qc.d(a,s+4)},function(a,b,d,f){a/=32;b/=32;d/=64;f/=64;b>=f+2-c.height()&&b<=f+2.5&&Math.abs(d-a)<=Math.max(c.width(),ic.width())/2?(g||O(),g=!0,t=n=p=l=k=h=!1):b>=f-3-e.height()&&b<=f-3&&Math.abs(d+2-a)<=e.width()/2?(h||O(),g=!1,h=!0,t=n=p=l=k=!1):b>=f-2-kc.height()&&b<=f-2&&Math.abs(d+2-a)<=kc.width()/2?(k||O(),h=g=!1,k=!0,t=n=p=l=!1):b>=f-1-mc.height()&&b<=f-1&&Math.abs(d+2-a)<=mc.width()/2?(l||O(),k=h=g=!1,l=!0,t=n=p=!1):b>=f-
oc.height()&&b<=f&&Math.abs(d+2-a)<=oc.width()/2?(p||O(),l=k=h=g=!1,p=!0,t=n=!1):b>=f+4-qc.height()&&b<=f+4&&Math.abs(d-a)<=qc.width()/2?(n||O(),p=l=k=h=g=!1,n=!0,t=!1):(t||O(),n=p=l=k=h=g=!1,t=!0);return!0},function(){switch(!0){case g:I({HUD:{N:"cc",D:"name"}});break;case h:I({HUD:{N:"cc",D:"gender"}});break;case k:I({HUD:{N:"cc",D:"skin"}});break;case l:I({HUD:{N:"cc",D:"shirt"}});break;case p:I({HUD:{N:"cc",D:"pants"}});break;case n:I({HUD:{N:"cc",D:"accept"}})}return!0},function(){return!0})};
U.examine=function(a){var b=new T(a.N,"#ccc",!0,!0),c=new T(a.E,"#ccc",!1,!0),d=[];(a.I||[]).forEach(function(a){var b=[];d.push(b);a.forEach(function(a){b.push(new T(a[0],a[1],!1,!0))})});return new gc("examine",function(a,f){L.fillStyle="rgba(0,0,0,.7)";L.fillRect(0,0,a,f);var g=a/2/32-6,h=f/2/32-6;b.d(g,h);h+=0.5;c.d(g,h);h+=1;d.forEach(function(b){g=a/2/32-6;b.forEach(function(a){a.d(g,h);g+=a.width()});h+=0.5})},function(){return!0},function(){V("examine");return!0},function(){V("examine");return!0})};
U.forge=function(a){var b=[];a.O.forEach(function(a){X.forEach(function(c){c.id==a.i&&b.push(c.e())})});b.forEach(function(a,b){a.x=a.L=b%8;a.y=a.M=Math.floor(b/8)});var c=[];a.C.forEach(function(a){c.push(a.map(bc))});var d=new Date(a.T),e=null,f=new $b("forge","no","",0,0,64,32),g=new T("Forge","#ccc",!0,!0),h=new T("Inventory","#aaa",!1,!0),k=null,l=null;return new gc("forge",function(a,e){L.fillStyle="rgba(0,0,0,.7)";L.fillRect(a/2-240,e/2-160,512,352);L.fillStyle="#000";L.beginPath();L.moveTo(a/
2-240,e/2-160);L.lineTo(a/2+272,e/2-160);L.lineTo(a/2+224,e/2-192);L.lineTo(a/2-232,e/2-192);L.closePath();L.fill();var t=a/32/2,w=e/32/2;f.d(t-7,w-5);g.d(t-5.5,w-5.1);var s=d-new Date;0<s&&(O(s%1E3),(new T(Math.floor(s/1E3)+"s remaining","#fff",!1,!0)).d(t-7.5,w-4.6));c.forEach(function(a,b){a.forEach(function(a){a.d(t-7.5+b%8,w-3.5+Math.floor(b/8))})});h.d(t+0.5,w-4.6);b.forEach(function(a){a.d(t+0.5,w-3.5)});k&&k.d(t+0.5,w+5);l&&l.d(t+0.5,w+5.5)},function(c,d,f,g){c=c/32-f/32/2;d=d/32-g/32/2;if(8<
Math.abs(c-0.5)||6<Math.abs(d))return!1;c=Math.floor(c-0.5);d=Math.floor(d+4.5);0<=c&&8>c&&0<=d&&c+8*d<b.length?(f=b[c+8*d],f!=e&&(e=f,k=new T(f.name,"#ccc",!1,!0),l=new T("quality "+a.O[c+8*d].q,"#aaa",!1,!0),O())):e&&(l=k=e=null,O());return!0},function(a,b,c,d){if(8<Math.abs(a/32-c/32/2-0.5)||6<Math.abs(b/32-d/32/2))return V("forge"),!1;e&&I({HUD:{N:"forge",D:{A:"a",I:e.id}}});return!0},function(a,b,c,d){return 8<Math.abs(a/32-c/32/2-0.5)||6<Math.abs(b/32-d/32/2)?(V("forge"),!1):!0})};
U.inv=function(){var a=-1,b=-1,c=-Infinity,d=-Infinity;return new gc("inv",function(e,f){var g=Math.ceil(X.length/8),h=e/32-8.1,k=f/32-0.1-g;L.fillStyle="rgba(0,0,0,.7)";L.fillRect(e-262.4,f-32*(g+1.2),262.4,32*(g+0.2));X.forEach(function(a){a.d(h,k)});if(0<=a&&8>a&&0<=b&&b<g&&a+8*b<X.length){var g=new T(X[a+8*b].name,"#fff",!1,!0),l=32*(g.width()+0.2),p=32*(g.height()+0.2);L.fillStyle="rgba(0,0,0,.7)";L.fillRect(c-l,d,l,p);g.d(c/32-g.width()-0.1,d/32+g.height()+0.1)}},function(e,f,g,h){var k=Math.ceil(X.length/
8);c=e;d=f;a=Math.floor((e-g)/32+8.1);b=Math.floor((f-h)/32+k+1.1);O();if(b>k||-1>a)V("inv");else return!0},function(){return!0},function(){var c=Math.ceil(X.length/8);0<=a&&8>a&&0<=b&&b<c&&a+8*b<X.length&&Zb("menu",X[a+8*b])})};
U.menu=function(a){var b=-1,c=null,d=null,e=!1,f=[],g=[],h=a.l;h.forEach(function(b){b=b+" "+a.name;f.push(new T(b,"#aaa",!1,!0));g.push(new T(b,"#fff",!1,!0))});var k=-2;f.forEach(function(a){k=Math.max(a.width(),k)});var k=32*(k+0.2),l=Math.floor(32*(h.length/2+0.2));return new gc("menu",function(a,n){e&&(e=!1,c+k>a&&(c=a-k),d+l>n&&(d=n-l));L.fillStyle="rgba(0,0,0,.7)";L.fillRect(c||0,d||0,k,l);-1!=b&&(L.fillStyle="#000",L.fillRect(c||0,(d||0)+32*(b/2+0.1),k,16));for(var t=c/32+0.1,w=d/32+0.5,s=
0;s<h.length;s++)s==b?g[s].d(t,w+s/2):f[s].d(t,w+s/2)},function(a,f){null===c&&(c=a,d=f,e=!0,O());a>=c&&a<c+k&&f>=Math.floor(d+3.2)&&f<Math.floor(d+32*(h.length/2+0.1))?(b=Math.floor(2*(f-d)/32-0.1),b>=h.length&&(b=-1)):(b=-1,(a<c-32||a>c+k+32||f<Math.floor(d-32)||f>d+l+32)&&V("menu"));O();return!0},function(){-1!=b&&I({Interact:{ID:a.id,X:a.x,Y:a.y,Action:h[b]}});V("menu");return!0},function(){V("menu");return!0})};
U.menu2=function(a){var b=-1,c=null,d=null,e=!1,f=[],g=[];a.forEach(function(a){a=a.name;f.push(new T(a,"#aaa",!1,!0));g.push(new T(a,"#fff",!1,!0))});var h=-2;f.forEach(function(a){h=Math.max(a.width(),h)});var h=32*(h+0.2),k=Math.floor(32*(a.length/2+0.2));return new gc("menu2",function(l,p){e&&(e=!1,c+h>l&&(c=l-h),d+k>p&&(d=p-k));L.fillStyle="rgba(0,0,0,.7)";L.fillRect(c||0,d||0,h,k);-1!=b&&(L.fillStyle="#000",L.fillRect(c||0,(d||0)+32*(b/2+0.1),h,16));for(var n=c/32+0.1,t=d/32+0.5,w=0;w<a.length;w++)w==
b?g[w].d(n,t+w/2):f[w].d(n,t+w/2)},function(f,g){null===c&&(c=f,d=g,e=!0,O());f>=c&&f<c+h&&g>=Math.floor(d+3.2)&&g<d+Math.floor(32*(a.length/2+0.1))?(b=Math.floor(2*(g-d)/32-0.1),b>=a.length&&(b=-1)):(b=-1,(f<c-32||f>c+h+32||g<d-32||g>d+k+32)&&V("menu2"));O();return!0},function(){-1!=b&&Zb("menu",a[b]);V("menu2");return!0},function(){V("menu2");return!0})};J.ClientHash=function(a){void 0!==Jc?Jc!=a&&location.reload(!0):Jc=a};
