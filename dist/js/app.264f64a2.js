(function(t){function e(e){for(var a,s,l=e[0],o=e[1],c=e[2],p=0,d=[];p<l.length;p++)s=l[p],Object.prototype.hasOwnProperty.call(i,s)&&i[s]&&d.push(i[s][0]),i[s]=0;for(a in o)Object.prototype.hasOwnProperty.call(o,a)&&(t[a]=o[a]);u&&u(e);while(d.length)d.shift()();return n.push.apply(n,c||[]),r()}function r(){for(var t,e=0;e<n.length;e++){for(var r=n[e],a=!0,l=1;l<r.length;l++){var o=r[l];0!==i[o]&&(a=!1)}a&&(n.splice(e--,1),t=s(s.s=r[0]))}return t}var a={},i={app:0},n=[];function s(e){if(a[e])return a[e].exports;var r=a[e]={i:e,l:!1,exports:{}};return t[e].call(r.exports,r,r.exports,s),r.l=!0,r.exports}s.m=t,s.c=a,s.d=function(t,e,r){s.o(t,e)||Object.defineProperty(t,e,{enumerable:!0,get:r})},s.r=function(t){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(t,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(t,"__esModule",{value:!0})},s.t=function(t,e){if(1&e&&(t=s(t)),8&e)return t;if(4&e&&"object"===typeof t&&t&&t.__esModule)return t;var r=Object.create(null);if(s.r(r),Object.defineProperty(r,"default",{enumerable:!0,value:t}),2&e&&"string"!=typeof t)for(var a in t)s.d(r,a,function(e){return t[e]}.bind(null,a));return r},s.n=function(t){var e=t&&t.__esModule?function(){return t["default"]}:function(){return t};return s.d(e,"a",e),e},s.o=function(t,e){return Object.prototype.hasOwnProperty.call(t,e)},s.p="/ui/";var l=window["webpackJsonp"]=window["webpackJsonp"]||[],o=l.push.bind(l);l.push=e,l=l.slice();for(var c=0;c<l.length;c++)e(l[c]);var u=o;n.push([0,"chunk-vendors"]),r()})({0:function(t,e,r){t.exports=r("cd49")},4080:function(t,e,r){},cd49:function(t,e,r){"use strict";r.r(e);r("e260"),r("e6cf"),r("cca6"),r("a79d"),r("caad"),r("2532");var a,i,n,s=r("2b0e"),l=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-app",{attrs:{id:"inspire"}},[r("v-navigation-drawer",{attrs:{app:"",width:"240"}},[r("v-sheet",{staticClass:"pa-4",attrs:{color:"light-blue lighten-3"}},[r("v-img",{attrs:{src:"https://raw.githubusercontent.com/falcosecurity/falcosidekick/master/imgs/falcosidekick_color.png",contain:"",width:"128"}}),r("v-divider",{staticClass:"my-3"}),r("h1",{staticClass:"text-h5"},[t._v("Falcosidekick UI")]),r("h2",{staticClass:"text-subtitle-1"},[t._v("Latest Events from Falco")]),r("v-chip",{staticClass:"my-2 mr-2 pr-0 py-0",staticStyle:{border:"1px solid #000"},attrs:{color:"primary",label:""}},[t._v(" Retention "),r("div",{staticClass:"white black--text rounded ml-2 px-2 py-2 text-body-2"},[t._v(t._s(t.retention))])]),r("v-divider",{staticClass:"my-3"}),r("h2",{staticClass:"text-subtitle-1"},[t._v("Outputs")]),t._l(t.outputs,(function(e){return r("v-chip",{key:e,staticClass:"my-2 mr-2",attrs:{color:"primary",label:""}},[t._v(" "+t._s(e)+" ")])}))],2),r("v-divider"),r("v-list",{staticClass:"py-0"},[r("v-list-item",{attrs:{to:"/",color:"light-blue"}},[r("v-list-item-content",[r("v-list-item-title",[t._v("Dashboard")])],1)],1),r("v-divider"),r("v-list-item",{attrs:{to:"/events",color:"light-blue"}},[r("v-list-item-content",[r("v-list-item-title",[t._v("Events")])],1)],1)],1),r("v-divider")],1),r("v-main",{staticClass:"grey lighten-3"},[r("router-view")],1)],1)},o=[],c=r("2f62"),u=r("ade3"),p=r("1da1"),d=(r("96cf"),r("bc3a")),h=r.n(d),v=r("3835");r("d81d"),r("13d5"),r("4fad"),r("d3b7"),r("25f0"),r("c1f9");(function(t){t["NONE"]="none",t["DEBUG"]="debug",t["INFORMATIONAL"]="informational",t["NOTICE"]="notice",t["WARNING"]="warning",t["ERROR"]="error",t["CRITICAL"]="critical",t["ALERT"]="alert",t["EMERGENCY"]="emergency"})(a||(a={}));var m,f=function(t){return t.map((function(t){var e=Object.entries(t.output_fields||{}).reduce((function(t,e){var r=Object(v["a"])(e,2),a=r[0],i=r[1];return t[a]=new String(i).toString(),t}),{}),r=Object.entries(e);return r.sort((function(t,e){var r=Object(v["a"])(t,2),a=r[0],i=r[1],n=Object(v["a"])(e,2),s=n[0],l=n[1];return a.length+i.length-(s.length+l.length)})),{output:t.output,priority:t.priority.toLowerCase(),rule:t.rule,time:new Date(t.time),outputFields:Object.fromEntries(r),filterTime:1e3*Math.floor(new Date(t.time).getTime()/1e3)}})).reverse()},b=function(){var t;return t={},Object(u["a"])(t,a.EMERGENCY,0),Object(u["a"])(t,a.ALERT,0),Object(u["a"])(t,a.CRITICAL,0),Object(u["a"])(t,a.ERROR,0),Object(u["a"])(t,a.WARNING,0),Object(u["a"])(t,a.NOTICE,0),Object(u["a"])(t,a.INFORMATIONAL,0),Object(u["a"])(t,a.DEBUG,0),Object(u["a"])(t,a.NONE,0),t},y=(i={},Object(u["a"])(i,a.NONE,"grey"),Object(u["a"])(i,a.DEBUG,"light-blue lighten-1"),Object(u["a"])(i,a.INFORMATIONAL,"light-blue"),Object(u["a"])(i,a.NOTICE,"primary"),Object(u["a"])(i,a.WARNING,"warning"),Object(u["a"])(i,a.ERROR,"error"),Object(u["a"])(i,a.CRITICAL,"red darken-1"),Object(u["a"])(i,a.ALERT,"red darken-2"),Object(u["a"])(i,a.EMERGENCY,"red darken-3"),i),O=(n={total:"#757575"},Object(u["a"])(n,a.EMERGENCY,"#C62828"),Object(u["a"])(n,a.ALERT,"#D32F2F"),Object(u["a"])(n,a.CRITICAL,"#E53935"),Object(u["a"])(n,a.ERROR,"#FF5252"),Object(u["a"])(n,a.WARNING,"#FB8C00"),Object(u["a"])(n,a.NOTICE,"#1976D2"),Object(u["a"])(n,a.INFORMATIONAL,"#03A9F4"),Object(u["a"])(n,a.DEBUG,"#29B6F6"),Object(u["a"])(n,a.NONE,"#555"),n),g=function(t){return y[t]},x=!0,C=h.a.create({baseURL:x?"//".concat(window.location.host,"/"):"http://localhost:2802"}),j={event:function(){return Object(p["a"])(regeneratorRuntime.mark((function t(){var e,r;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.next=2,C.get("events");case 2:return e=t.sent,r=e.data,t.abrupt("return",{events:f(r.events||[]),retention:r.retention,stats:r.stats,outputs:r.outputs});case 5:case"end":return t.stop()}}),t)})))()}},E=r("f7fe"),_=r.n(E);s["a"].use(c["a"]);var w="FETCH_EVENTS",T="RECEIVE_EVENTS",S="SET_STATE",k=new c["a"].Store({state:{socket:{isConnected:!1},events:[],outputs:[],retention:0},mutations:Object(u["a"])({},S,(function(t,e){t.events=e.events||[],t.outputs=e.outputs||[],t.retention=e.retention})),actions:(m={},Object(u["a"])(m,w,(function(t){var e=t.commit;j.event().then((function(t){return e(S,t)}))})),Object(u["a"])(m,T,_()((function(t,e){var r,i=t.commit,n={events:f(e.events||[]),outputs:e.outputs,retention:e.retention};e.stats&&(n.stats=(r={},Object(u["a"])(r,a.EMERGENCY,e.stats[a.EMERGENCY]),Object(u["a"])(r,a.ALERT,e.stats[a.ALERT]),Object(u["a"])(r,a.CRITICAL,e.stats[a.CRITICAL]),Object(u["a"])(r,a.ERROR,e.stats[a.ERROR]),Object(u["a"])(r,a.WARNING,e.stats[a.WARNING]),Object(u["a"])(r,a.NOTICE,e.stats[a.NOTICE]),Object(u["a"])(r,a.INFORMATIONAL,e.stats[a.INFORMATIONAL]),Object(u["a"])(r,a.DEBUG,e.stats[a.DEBUG]),Object(u["a"])(r,a.NONE,e.stats[a.NONE]),r));i(S,n)}),1e3)),m)}),R=s["a"].extend({name:"App",computed:Object(c["b"])(["retention","outputs"]),created:function(){this.$store.dispatch(w)}}),A=R,V=r("2877"),N=r("6544"),P=r.n(N),$=r("7496"),I=r("cc20"),L=r("ce7e"),M=r("adda"),D=r("8860"),F=r("da13"),q=r("5d23"),G=r("f6c4"),z=r("f774"),B=r("8dd9"),H=Object(V["a"])(A,l,o,!1,null,null,null),W=H.exports;P()(H,{VApp:$["a"],VChip:I["a"],VDivider:L["a"],VImg:M["a"],VList:D["a"],VListItem:F["a"],VListItemContent:q["a"],VListItemTitle:q["b"],VMain:G["a"],VNavigationDrawer:z["a"],VSheet:B["a"]});var U=r("8c4f"),Y=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-container",{staticClass:"px-3 py-2",attrs:{fluid:""}},[r("v-row",{staticClass:"pb-0"},[r("v-col",{attrs:{cols:"12"}},[r("v-card",{attrs:{elevation:"2"}},[r("v-toolbar",{attrs:{flat:""}},[r("v-spacer"),r("rule-autocomplete",{staticClass:"ml-3",staticStyle:{"max-width":"500px"},attrs:{events:t.events},model:{value:t.rules,callback:function(e){t.rules=e},expression:"rules"}}),r("priority-autocomplete",{staticClass:"ml-3",staticStyle:{"max-width":"350px"},model:{value:t.priorities,callback:function(e){t.priorities=e},expression:"priorities"}}),r("time-range-select",{staticClass:"ml-3",staticStyle:{"max-width":"300px"},model:{value:t.time,callback:function(e){t.time=e},expression:"time"}})],1)],1)],1)],1),r("v-row",{staticClass:"pb-0 mt-0"},[r("v-col",{attrs:{cols:"4"}},[r("priority-chart",{attrs:{events:t.filtered}})],1),r("v-col",{attrs:{cols:"8"}},[r("rule-chart",{attrs:{events:t.filtered}})],1)],1),r("v-row",{staticClass:"pb-0 mt-0"},[r("v-col",{attrs:{cols:"12"}},[r("rule-time-chart",{key:t.time,attrs:{events:t.filtered,time:t.time}})],1)],1),r("v-row",{staticClass:"pb-0 mt-0"},[r("v-col",{attrs:{cols:"12"}},[r("priority-time-chart",{key:t.time,attrs:{events:t.filtered,time:t.time}})],1)],1)],1)},J=[],K=r("5530"),Q=(r("b0c0"),r("4de4"),function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-autocomplete",t._b({attrs:{dense:"",multiple:"",items:t.options,outlined:"","hide-details":"",clearable:"",name:"Priorities",label:"Priorities",value:t.value},on:{input:t.input},scopedSlots:t._u([{key:"selection",fn:function(e){var a=e.item,i=e.index;return[i<=1?r("v-chip",{attrs:{small:"",label:"",outlined:""}},[r("span",[t._v(t._s(a))])]):t._e(),2===i?r("span",{staticClass:"grey--text caption"},[t._v(" (+"+t._s(t.selected.length-2)+" others) ")]):t._e()]}}])},"v-autocomplete",t.$attrs,!1))}),X=[],Z=(r("b64b"),_()((function(t){t()}),800)),tt=s["a"].extend({props:{value:{type:Array,default:function(){return[]}}},data:function(){return{options:Object.keys(b()),selected:this.value}},methods:{input:function(t){var e=this;this.selected=t,Z((function(){e.$emit("input",t)}))}}}),et=tt,rt=r("c6a6"),at=Object(V["a"])(et,Q,X,!1,null,null,null),it=at.exports;P()(at,{VAutocomplete:rt["a"],VChip:I["a"]});var nt=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-card",{staticStyle:{height:"100%"}},[r("v-card-title",[t._v(" Priorities ")]),r("v-card-text",[r("wait",[r("apexchart",{attrs:{type:"donut",height:"350",options:t.pie.chartOptions,series:t.pie.series}})],1)],1),r("event-overlay",{model:{value:t.selectedEvents,callback:function(e){t.selectedEvents=e},expression:"selectedEvents"}})],1)},st=[],lt=r("2909"),ot=(r("a9e3"),r("159b"),r("07ac"),function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-dialog",{model:{value:t.dialog,callback:function(e){t.dialog=e},expression:"dialog"}},[r("v-card",[r("v-card-title",[t._v(" Events "),r("v-spacer"),r("v-text-field",{staticClass:"ml-3",staticStyle:{"max-width":"300px"},attrs:{outlined:"","hide-details":"",dense:"",label:"Search",clearable:""},model:{value:t.search,callback:function(e){t.search=e},expression:"search"}}),r("v-btn",{staticClass:"ml-4",attrs:{icon:""},on:{click:t.close}},[r("v-icon",[t._v("mdi-close")])],1)],1),r("v-divider"),t.value?r("event-table",{attrs:{events:t.value,search:t.search},on:{"update:search":function(e){t.search=e}}}):t._e(),r("v-divider"),r("v-card-actions",[r("v-spacer"),r("v-btn",{attrs:{color:"primary",text:""},on:{click:t.close}},[t._v(" Close ")])],1)],1)],1)}),ct=[],ut=(r("841c"),r("ac1f"),function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-data-table",{attrs:{items:t.events,headers:t.headers,"items-per-page":-1,"hide-default-footer":"",search:t.search,"custom-filter":t.defaultFilter,"disable-sort":""},scopedSlots:t._u([{key:"item",fn:function(e){var a=e.item;return[r("tr",[r("td",[r("span",{domProps:{innerHTML:t._s(a.time.toISOString().split("T").join(" ").slice(0,-5))}})]),r("td",[r("v-chip",{staticClass:"mb-3 white--text",attrs:{label:"",color:t.color(a.priority)},on:{click:function(e){return t.$emit("update:search",a.priority)}}},[r("span",{domProps:{innerHTML:t._s(t.highlightMatches(a.priority))}})]),r("br"),r("v-chip",{staticClass:"black--text",attrs:{label:"",color:"light-blue lighten-5"},domProps:{innerHTML:t._s(t.highlightMatches(a.rule))},on:{click:function(e){return t.$emit("update:search",a.rule)}}}),r("br")],1),r("td",{staticClass:"py-3",staticStyle:{"font-size":"0.8rem"}},[r("div",{domProps:{innerHTML:t._s(t.highlightMatches(a.output))}}),r("div",{staticClass:"mt-3"},[t._l(a.outputFields,(function(e,a){return[e.length<=100?r("output-field-chip",{key:a,attrs:{dark:"",value:t.highlightMatches(e),label:a},on:{click:function(r){return t.$emit("update:search",e)}}}):r("output-field-card",{key:a,attrs:{value:t.highlightMatches(e),label:a},on:{click:function(r){return t.$emit("update:search",e)}}})]}))],2)])])]}}])})}),pt=[],dt=(r("5319"),r("4d63"),function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-card",t._g(t._b({staticClass:"primary mt-2",staticStyle:{border:"1px solid"},attrs:{flat:""}},"v-card",t.$attrs,!1),t.$listeners),[r("v-system-bar",{attrs:{color:"primary white--text"}},[t._v(t._s(t.label))]),r("v-card-text",{staticClass:"pa-2 white black--text",staticStyle:{"font-size":"0.8rem","line-height":"1.1rem"}},[r("span",{domProps:{innerHTML:t._s(t.value)}})])],1)}),ht=[],vt=s["a"].extend({name:"OutputFieldCard",inheritAttrs:!1,props:{label:{type:String,required:!0},value:{type:String,required:!0}}}),mt=vt,ft=r("b0af"),bt=r("99d9"),yt=r("afd9"),Ot=Object(V["a"])(mt,dt,ht,!1,null,null,null),gt=Ot.exports;P()(Ot,{VCard:ft["a"],VCardText:bt["b"],VSystemBar:yt["a"]});var xt=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-chip",t._g({staticClass:"my-1 mr-2 pr-0 py-0",staticStyle:{border:"1px solid #000"},attrs:{dark:"",label:"",small:"",color:"primary"}},t.$listeners),[t._v(" "+t._s(t.label)+" "),r("div",{staticClass:"white black--text rounded ml-2 px-2 py-2 text-body-2",staticStyle:{"font-size":"0.7rem!important","line-height":"0.9rem!important"}},[r("span",{domProps:{innerHTML:t._s(t.value)}})])])},Ct=[],jt=s["a"].extend({name:"OutputFieldChip",props:{label:{type:String,required:!0},value:{type:String,required:!0}}}),Et=jt,_t=Object(V["a"])(Et,xt,Ct,!1,null,null,null),wt=_t.exports;P()(_t,{VChip:I["a"]});var Tt=s["a"].extend({components:{OutputFieldChip:wt,OutputFieldCard:gt},name:"EventTable",props:{events:{type:Array,default:function(){return[]}},search:{type:String}},data:function(){return{headers:[{text:"Time",value:"time"},{text:"Rule / Priority",value:"rule"},{text:"Output",value:"output",width:"60%"}]}},methods:{color:function(t){return g(t)},highlightMatches:function(t){if(!this.search)return t;var e=t.toLowerCase().includes(this.search.toLowerCase());return e?t.replace(new RegExp(this.search,"ig"),(function(t){return"<strong>".concat(t,"</strong>")})):t},defaultFilter:function(t,e,r){e=(e||"").toLocaleLowerCase();var a=Object.values(r.outputFields).some((function(t){return t.toLowerCase().includes(e)}));if(a)return!0;var i=null!=r.priority&&null!=e&&-1!==r.priority.toLocaleLowerCase().indexOf(e);return!!i||null!=t&&null!=e&&"boolean"!==typeof t&&-1!==t.toString().toLocaleLowerCase().indexOf(e)}}}),St=Tt,kt=(r("fedc"),r("8fea")),Rt=Object(V["a"])(St,ut,pt,!1,null,"2402e04a",null),At=Rt.exports;P()(Rt,{VChip:I["a"],VDataTable:kt["a"]});var Vt=s["a"].extend({components:{EventTable:At},props:{value:{type:Array,default:function(){return[]}}},data:function(){return{search:"",dialog:!1}},watch:{value:function(t){var e=this;0!==t.length&&(this.search="",setTimeout((function(){e.dialog=!0}),200))}},methods:{close:function(){var t=this;this.dialog=!1,setTimeout((function(){t.$emit("input",[])}),200)}}}),Nt=Vt,Pt=r("8336"),$t=r("169a"),It=r("132d"),Lt=r("2fa4"),Mt=r("8654"),Dt=Object(V["a"])(Nt,ot,ct,!1,null,null,null),Ft=Dt.exports;P()(Dt,{VBtn:Pt["a"],VCard:ft["a"],VCardActions:bt["a"],VCardTitle:bt["c"],VDialog:$t["a"],VDivider:L["a"],VIcon:It["a"],VSpacer:Lt["a"],VTextField:Mt["a"]});var qt=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("div",[t.ready?t._t("default"):t._e()],2)},Gt=[],zt=s["a"].extend({data:function(){return{ready:!1}},created:function(){var t=this;setTimeout((function(){t.ready=!0}),100)}}),Bt=zt,Ht=Object(V["a"])(Bt,qt,Gt,!1,null,null,null),Wt=Ht.exports,Ut=s["a"].extend({components:{Wait:Wt,EventOverlay:Ft},name:"PieChart",props:{events:{type:Array,required:!0},time:{type:Number}},data:function(){return{selectedEvents:[],stats:b(),eventsPerPriority:{}}},watch:{events:{immediate:!0,handler:function(){var t=b(),e={};this.events.forEach((function(r){t[r.priority]+=1,!1===Object.prototype.hasOwnProperty.call(e,r.priority)&&(e[r.priority]=[]),e[r.priority].push(r)})),this.stats=t,this.eventsPerPriority=e}}},computed:{pie:function(){var t=this,e=Object.keys(this.stats);return{series:Object.values(this.stats),chartOptions:{chart:{height:350,type:"donut",selection:{enabled:!1},events:{dataPointSelection:function(r,a,i){t.selectedEvents=Object(lt["a"])(t.eventsPerPriority[e[i.dataPointIndex]])}}},states:{active:{allowMultipleDataPointsSelection:!1,filter:{type:"none"}}},dataLabels:{dropShadow:{enabled:!0,top:0,left:0,blur:1,color:"#000",opacity:1}},plotOptions:{pie:{expandOnClick:!1,donut:{labels:{show:!0,total:{showAlways:!0,show:!0}}}}},labels:e,colors:Object.keys(this.stats).map((function(t){return O[t]}))}}}}}),Yt=Ut,Jt=Object(V["a"])(Yt,nt,st,!1,null,null,null),Kt=Jt.exports;P()(Jt,{VCard:ft["a"],VCardText:bt["b"],VCardTitle:bt["c"]});var Qt=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-card",{staticStyle:{height:"100%"}},[r("v-card-title",[t._v(" Priority Timeline ")]),r("v-card-text",[r("wait",[r("apexchart",{attrs:{type:"bar",height:t.height,options:t.chart,series:t.series}})],1)],1),r("event-overlay",{model:{value:t.selectedEvents,callback:function(e){t.selectedEvents=e},expression:"selectedEvents"}})],1)},Xt=[],Zt=(r("a630"),r("3ca3"),r("ddb0"),r("fb6a"),r("1276"),function(t,e){var r=(new Date).getTime()-t;return r-r%e}),te=function(t){var e=arguments.length>1&&void 0!==arguments[1]?arguments[1]:12,r=t||3e5,a=r/e;return{range:r,steps:e,start:Zt(r,a),stepSize:a}},ee=s["a"].extend({components:{Wait:Wt,EventOverlay:Ft},name:"PriorityTimeChart",props:{events:{type:Array,required:!0},time:{type:Number}},data:function(){var t=te(this.time,15),e=t.range,r=t.steps,a=t.start,i=t.stepSize;return{range:e,steps:r,start:a,stepSize:i,height:400,series:[],selectedEvents:[],eventsPerBucket:{},categories:[]}},watch:{events:{immediate:!0,handler:function(t){var e=this,r=Zt(this.range,this.stepSize),a=Array.from(Array(this.steps).keys()).reduce((function(t,a){var i=r+(a+1)*e.stepSize,n=new Date(i).toISOString().replace("T"," ").slice(0,-5);return e.range>=36e5&&(n="".concat(n.slice(0,-2),"00")),t.push({time:i,label:n}),t}),[]),i={},n={};Object.keys(b()).forEach((function(t){n[t]={name:t,data:Object(lt["a"])(Array(e.steps).keys()).map((function(){return 0}))},i[t]=Object(lt["a"])(Array(e.steps).keys()).map((function(){return[]}))})),t.forEach((function(t){if(!(t.filterTime<e.start)){for(var r in a)if(a[+r+1]&&t.filterTime<a[+r+1].time)return n[t.priority].data[r]+=1,void i[t.priority][r].push(t);n[t.priority].data[e.steps-1]+=1,i[t.priority][e.steps-1].push(t)}})),this.series=Object.values(n).reverse(),this.eventsPerBucket=i,this.categories=a,this.start=r}}},computed:{chart:function(){var t=this,e=this.series.map((function(t){return t.name}));return{chart:{type:"bar",height:this.height,stacked:!0,toolbar:{show:!1},events:{dataPointSelection:function(r,a,i){var n=e[i.seriesIndex];t.selectedEvents=Object(lt["a"])(t.eventsPerBucket[n][i.dataPointIndex])}}},colors:this.series.map((function(t){return O[t.name]})),plotOptions:{bar:{horizontal:!1}},states:{active:{allowMultipleDataPointsSelection:!1,filter:{type:"none"}}},xaxis:{categories:this.categories.map((function(t){var e=t.label;return e.split(" ")}))},legend:{position:"right",offsetY:40},fill:{opacity:1}}}}}),re=ee,ae=Object(V["a"])(re,Qt,Xt,!1,null,null,null),ie=ae.exports;P()(ae,{VCard:ft["a"],VCardText:bt["b"],VCardTitle:bt["c"]});var ne=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-autocomplete",t._b({attrs:{dense:"",multiple:"",items:t.options,outlined:"","hide-details":"",name:"Rules",label:"Rules",clearable:"",value:t.value},on:{input:t.input},scopedSlots:t._u([{key:"selection",fn:function(e){var a=e.item,i=e.index;return[0===i?r("v-chip",{attrs:{small:"",label:"",outlined:""}},[r("span",[t._v(t._s(a))])]):t._e(),1===i?r("span",{staticClass:"grey--text caption"},[t._v(" (+"+t._s(t.value.length-1)+" others) ")]):t._e()]}}])},"v-autocomplete",t.$attrs,!1))},se=[],le=_()((function(t){t()}),800),oe=s["a"].extend({props:{events:{type:Array,default:function(){return[]}},value:{type:Array,default:function(){return[]}}},data:function(){return{selected:[]}},computed:{options:function(){return Object.keys(this.events.reduce((function(t,e){return t[e.rule]=null,t}),{})).sort()}},methods:{input:function(t){var e=this;this.selected=t,le((function(){e.$emit("input",t)}))}}}),ce=oe,ue=Object(V["a"])(ce,ne,se,!1,null,null,null),pe=ue.exports;P()(ue,{VAutocomplete:rt["a"],VChip:I["a"]});var de=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-card",{staticStyle:{height:"100%"}},[r("v-card-title",[t._v(" Rules ")]),r("v-card-text",[r("apexchart",{attrs:{type:"bar",height:"350",options:t.chart.chartOptions,series:t.chart.series}})],1),r("event-overlay",{model:{value:t.selectedEvents,callback:function(e){t.selectedEvents=e},expression:"selectedEvents"}})],1)},he=[],ve=(r("a15b"),function(t){return t.split(" ").reduce((function(t,e,r){var a=Math.floor(r/3);return t[a]||(t[a]=[]),t[a].push(e),t}),[]).map((function(t){return t.join(" ")}))}),me=s["a"].extend({components:{EventOverlay:Ft},name:"RuleChart",props:{events:{type:Array,required:!0},time:{type:Number}},data:function(){return{selectedEvents:[],eventsPerRule:{},rules:{}}},watch:{events:{immediate:!0,handler:function(t){var e={},r={};t.forEach((function(t){return!1===Object.prototype.hasOwnProperty.call(e,t.rule)&&(e[t.rule]=0,r[t.rule]=[]),e[t.rule]+=1,r[t.rule].push(t),e})),this.rules=Object.keys(e).sort().reduce((function(t,r){return t[r]=e[r],t}),{}),this.eventsPerRule=r}}},computed:{chart:function(){var t=this,e=Object.keys(this.rules);return{series:[{name:"Counter",data:Object.values(this.rules)}],chartOptions:{chart:{height:350,type:"bar",toolbar:{show:!1},selection:{enabled:!1},zoom:{enabled:!1},events:{dataPointSelection:function(r,a,i){t.selectedEvents=Object(lt["a"])(t.eventsPerRule[e[i.dataPointIndex]])}}},dataLabels:{enabled:!1},states:{active:{allowMultipleDataPointsSelection:!1,filter:{type:"none"}}},xaxis:{labels:{rotate:-45,rotateAlways:!0,hideOverlappingLabels:!1,maxHeight:200,style:{colors:[],fontSize:"12px"}},categories:Object.keys(this.rules).map((function(t){return ve(t)})),tickPlacement:"on"}}}}}}),fe=me,be=Object(V["a"])(fe,de,he,!1,null,null,null),ye=be.exports;P()(be,{VCard:ft["a"],VCardText:bt["b"],VCardTitle:bt["c"]});var Oe=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-card",{staticStyle:{height:"100%"}},[r("v-card-title",[t._v(" Rule Timeline ")]),r("v-card-text",[r("wait",[r("apexchart",{attrs:{type:"bar",height:t.height,options:t.chart,series:t.series}})],1)],1),r("event-overlay",{model:{value:t.selectedEvents,callback:function(e){t.selectedEvents=e},expression:"selectedEvents"}})],1)},ge=[],xe=s["a"].extend({components:{Wait:Wt,EventOverlay:Ft},name:"RuleTimeChart",props:{events:{type:Array,required:!0},time:{type:Number}},data:function(){var t=te(this.time,15),e=t.range,r=t.steps,a=t.start,i=t.stepSize;return{range:e,steps:r,start:a,stepSize:i,height:400,series:[],selectedEvents:[],eventsPerBucket:{},categories:[]}},watch:{events:{immediate:!0,handler:function(t){var e=this,r=Zt(this.range,this.stepSize),a=Array.from(Array(this.steps).keys()).reduce((function(t,a){var i=r+(a+1)*e.stepSize,n=new Date(i).toISOString().replace("T"," ").slice(0,-5);return e.range>=36e5&&(n="".concat(n.slice(0,-2),"00")),t.push({time:i,label:n}),t}),[]),i={},n={};t.forEach((function(t){if(!1===Object.prototype.hasOwnProperty.call(n,t.rule)&&(n[t.rule]={name:t.rule,data:Object(lt["a"])(Array(e.steps).keys()).map((function(){return 0}))},i[t.rule]=Object(lt["a"])(Array(e.steps).keys()).map((function(){return[]}))),!(t.filterTime<r)){for(var s in a)if(a[+s+1]&&t.filterTime<a[+s+1].time)return n[t.rule].data[s]+=1,void i[t.rule][s].push(t);n[t.rule].data[e.steps-1]+=1,i[t.rule][e.steps-1].push(t)}}));var s=Object.keys(n).sort().reduce((function(t,e){return t[e]=n[e],t}),{});this.eventsPerBucket=i,this.categories=a,this.start=r,this.series=Object.values(s).filter((function(t){var e=t.data.reduce((function(t,e){return t+e}),0);return e>0}))}}},computed:{chart:function(){var t=this,e=this.series.map((function(t){return t.name}));return{chart:{type:"bar",height:this.height,stacked:!0,toolbar:{show:!1},zoom:{enabled:!1},events:{dataPointSelection:function(r,a,i){var n=e[i.seriesIndex];t.selectedEvents=Object(lt["a"])(t.eventsPerBucket[n][i.dataPointIndex])}}},plotOptions:{bar:{horizontal:!1}},states:{active:{allowMultipleDataPointsSelection:!1,filter:{type:"none"}}},xaxis:{categories:this.categories.map((function(t){var e=t.label;return e.split(" ")}))},legend:{position:"right",offsetY:40},fill:{opacity:1}}}}}),Ce=xe,je=Object(V["a"])(Ce,Oe,ge,!1,null,null,null),Ee=je.exports;P()(je,{VCard:ft["a"],VCardText:bt["b"],VCardTitle:bt["c"]});var _e=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-select",t._g(t._b({attrs:{outlined:"","hide-details":"",dense:"",label:"Timerange","item-value":"value","item-text":"label",items:t.options,clearable:""}},"v-select",t.$attrs,!1),t.$listeners))},we=[],Te=s["a"].extend({inheritAttrs:!1,data:function(){return{options:[{label:"last 15 minutes",value:9e5},{label:"last 30 minutes",value:18e5},{label:"last hour",value:36e5},{label:"last 6 hours",value:216e5},{label:"last 12 hours",value:432e5},{label:"last 24 hours",value:864e5},{label:"last 2 days",value:1728e5},{label:"last 7 days",value:6048e5}]}}}),Se=Te,ke=r("b974"),Re=Object(V["a"])(Se,_e,we,!1,null,null,null),Ae=Re.exports;P()(Re,{VSelect:ke["a"]});var Ve=s["a"].extend({components:{TimeRangeSelect:Ae,PriorityChart:Kt,RuleChart:ye,RuleTimeChart:Ee,PriorityTimeChart:ie,PriorityAutocomplete:it,RuleAutocomplete:pe},name:"Dashboard",computed:Object(K["a"])(Object(K["a"])({},Object(c["b"])(["events"])),{},{time:{get:function(){var t=this.$route.query.time;return t?+t:null},set:function(t){t?this.$router.push({name:this.$route.name||"",query:Object(K["a"])(Object(K["a"])({},this.$route.query),{},{time:t.toString()})}):this.$router.push({name:this.$route.name||"",query:Object(K["a"])(Object(K["a"])({},this.$route.query),{},{time:void 0})})}},priorities:{get:function(){var t=this.$route.query.priority;return t?Array.isArray(t)?t:[t]:[]},set:function(t){t?this.$router.push({name:this.$route.name||"",query:Object(K["a"])(Object(K["a"])({},this.$route.query),{},{priority:t})}):this.$router.push({name:this.$route.name||"",query:Object(K["a"])(Object(K["a"])({},this.$route.query),{},{priority:void 0})})}},rules:{get:function(){var t=this.$route.query.rules;return t?Array.isArray(t)?t:[t]:[]},set:function(t){t?this.$router.push({name:this.$route.name||"",query:Object(K["a"])(Object(K["a"])({},this.$route.query),{},{rules:t})}):this.$router.push({name:this.$route.name||"",query:Object(K["a"])(Object(K["a"])({},this.$route.query),{},{rules:void 0})})}},filtered:function(){var t=this,e=this.time,r=null;return"number"===typeof e&&e>0&&(r=(new Date).getTime()-e),0!==this.rules.length||0!==this.priorities.length||r?this.events.filter((function(e){return(!r||e.filterTime>=r)&&(!t.priorities.length||t.priorities.includes(e.priority))&&(!t.rules.length||t.rules.includes(e.rule))})):this.events}})}),Ne=Ve,Pe=r("62ad"),$e=r("a523"),Ie=r("0fd9"),Le=r("71d9"),Me=Object(V["a"])(Ne,Y,J,!1,null,null,null),De=Me.exports;P()(Me,{VCard:ft["a"],VCol:Pe["a"],VContainer:$e["a"],VRow:Ie["a"],VSpacer:Lt["a"],VToolbar:Le["a"]});var Fe=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-container",{staticClass:"px-3 py-2",attrs:{fluid:""}},[r("v-row",{staticClass:"pb-0"},[r("v-col",{attrs:{cols:"12"}},[r("v-card",{attrs:{elevation:"2"}},[r("v-card-text",t._l(t.stats,(function(e,a){return r("priority-counter-chip",{key:a,attrs:{dark:"",priority:a,count:e},on:{click:function(e){t.search=a}}},[t._v(" "+t._s(a)+" "),r("div",{staticClass:"white black--text rounded ml-2 px-2 py-2 text-body-2"},[t._v(t._s(e))])])})),1)],1)],1)],1),r("v-row",{staticClass:"pt-0 mt-0"},[r("v-col",{attrs:{cols:"12"}},[r("v-card",{attrs:{elevation:"2"}},[r("v-toolbar",{attrs:{flat:""}},[r("v-spacer"),r("time-range-select",{staticClass:"ml-3",staticStyle:{"max-width":"300px"},model:{value:t.time,callback:function(e){t.time=e},expression:"time"}}),r("v-text-field",{staticClass:"ml-3",staticStyle:{"max-width":"300px"},attrs:{outlined:"","hide-details":"",dense:"",label:"Search",clearable:""},model:{value:t.search,callback:function(e){t.search=e},expression:"search"}})],1),r("v-divider"),r("event-table",{attrs:{events:t.filtered,search:t.search},on:{"update:search":function(e){t.search=e}}})],1)],1)],1)],1)},qe=[],Ge=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-chip",t._g({staticClass:"my-0 mr-2 pr-0 py-0",staticStyle:{border:"1px solid #000"},attrs:{dark:"",color:t.color(t.priority),label:""}},t.$listeners),[t._v(" "+t._s(t.priority)+" "),r("div",{staticClass:"white black--text rounded ml-2 px-2 py-2 text-body-2"},[t._v(t._s(t.count))])])},ze=[],Be=s["a"].extend({name:"PriorityCountChip",props:{priority:{type:String,required:!0},count:{type:Number,required:!0}},methods:{color:function(t){return g(t)}}}),He=Be,We=Object(V["a"])(He,Ge,ze,!1,null,null,null),Ue=We.exports;P()(We,{VChip:I["a"]});var Ye=s["a"].extend({components:{EventTable:At,PriorityCounterChip:Ue,TimeRangeSelect:Ae},name:"Dashboard",data:function(){return{search:"",time:null}},computed:Object(K["a"])(Object(K["a"])({},Object(c["b"])(["events"])),{},{filtered:function(){var t=this.time;if("number"===typeof t&&t>0){var e=(new Date).getTime()-t;return this.events.filter((function(t){return t.filterTime>=e}))}return this.events},stats:function(){var t=Object(K["a"])({total:0},b());return this.events.forEach((function(e){t[e.priority]+=1,t.total+=1})),t}}),methods:{color:function(t){return g(t)}}}),Je=Ye,Ke=Object(V["a"])(Je,Fe,qe,!1,null,null,null),Qe=Ke.exports;P()(Ke,{VCard:ft["a"],VCardText:bt["b"],VCol:Pe["a"],VContainer:$e["a"],VDivider:L["a"],VRow:Ie["a"],VSpacer:Lt["a"],VTextField:Mt["a"],VToolbar:Le["a"]}),s["a"].use(U["a"]);var Xe=[{path:"/",name:"Dashboard",component:De},{path:"/events",name:"Events",component:Qe}],Ze=new U["a"]({mode:"hash",base:"/ui/",routes:Xe}),tr=Ze,er=r("f309");s["a"].use(er["a"]);var rr=new er["a"]({}),ar=r("b408"),ir=r.n(ar),nr=r("1321"),sr=r.n(nr);s["a"].use(sr.a),s["a"].component("apexchart",sr.a);var lr="ws://".concat(window.location.host,"/ws");s["a"].use(ir.a,lr,{format:"json",reconnection:!0,store:k,passToStoreHandler:function(t,e){t.includes("onmessage")&&this.store.dispatch(T,JSON.parse(e.data))}}),s["a"].config.productionTip=!1,new s["a"]({router:tr,store:k,vuetify:rr,render:function(t){return t(W)}}).$mount("#app")},fedc:function(t,e,r){"use strict";r("4080")}});
//# sourceMappingURL=app.264f64a2.js.map