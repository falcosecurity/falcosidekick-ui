(function(t){function e(e){for(var a,s,o=e[0],l=e[1],c=e[2],p=0,h=[];p<o.length;p++)s=o[p],Object.prototype.hasOwnProperty.call(i,s)&&i[s]&&h.push(i[s][0]),i[s]=0;for(a in l)Object.prototype.hasOwnProperty.call(l,a)&&(t[a]=l[a]);u&&u(e);while(h.length)h.shift()();return n.push.apply(n,c||[]),r()}function r(){for(var t,e=0;e<n.length;e++){for(var r=n[e],a=!0,o=1;o<r.length;o++){var l=r[o];0!==i[l]&&(a=!1)}a&&(n.splice(e--,1),t=s(s.s=r[0]))}return t}var a={},i={app:0},n=[];function s(e){if(a[e])return a[e].exports;var r=a[e]={i:e,l:!1,exports:{}};return t[e].call(r.exports,r,r.exports,s),r.l=!0,r.exports}s.m=t,s.c=a,s.d=function(t,e,r){s.o(t,e)||Object.defineProperty(t,e,{enumerable:!0,get:r})},s.r=function(t){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(t,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(t,"__esModule",{value:!0})},s.t=function(t,e){if(1&e&&(t=s(t)),8&e)return t;if(4&e&&"object"===typeof t&&t&&t.__esModule)return t;var r=Object.create(null);if(s.r(r),Object.defineProperty(r,"default",{enumerable:!0,value:t}),2&e&&"string"!=typeof t)for(var a in t)s.d(r,a,function(e){return t[e]}.bind(null,a));return r},s.n=function(t){var e=t&&t.__esModule?function(){return t["default"]}:function(){return t};return s.d(e,"a",e),e},s.o=function(t,e){return Object.prototype.hasOwnProperty.call(t,e)},s.p="/ui/";var o=window["webpackJsonp"]=window["webpackJsonp"]||[],l=o.push.bind(o);o.push=e,o=o.slice();for(var c=0;c<o.length;c++)e(o[c]);var u=l;n.push([0,"chunk-vendors"]),r()})({0:function(t,e,r){t.exports=r("cd49")},4080:function(t,e,r){},cd49:function(t,e,r){"use strict";r.r(e);r("e260"),r("e6cf"),r("cca6"),r("a79d"),r("caad"),r("2532");var a,i,n,s=r("2b0e"),o=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-app",{attrs:{id:"inspire"}},[r("v-navigation-drawer",{attrs:{app:"",width:"240"}},[r("v-sheet",{staticClass:"pa-4",attrs:{color:"light-blue lighten-3"}},[r("v-img",{attrs:{src:"https://raw.githubusercontent.com/falcosecurity/falcosidekick/master/imgs/falcosidekick_color.png",contain:"",width:"128"}}),r("v-divider",{staticClass:"my-3"}),r("h1",{staticClass:"text-h5"},[t._v("Falcosidekick UI")]),r("h2",{staticClass:"text-subtitle-1"},[t._v("Latest Events from Falco")]),r("v-chip",{staticClass:"my-2 mr-2 pr-0 py-0",staticStyle:{border:"1px solid #000"},attrs:{color:"primary",label:""}},[t._v(" Retention "),r("div",{staticClass:"white black--text rounded ml-2 px-2 py-2 text-body-2"},[t._v(t._s(t.retention))])]),r("v-divider",{staticClass:"my-3"}),r("h2",{staticClass:"text-subtitle-1"},[t._v("Outputs")]),t._l(t.outputs,(function(e){return r("v-chip",{key:e,staticClass:"my-2 mr-2",attrs:{color:"primary",label:""}},[t._v(" "+t._s(e)+" ")])}))],2),r("v-divider"),r("v-list",{staticClass:"py-0"},[r("v-list-item",{attrs:{to:"/",color:"light-blue"}},[r("v-list-item-content",[r("v-list-item-title",[t._v("Dashboard")])],1)],1),r("v-divider"),r("v-list-item",{attrs:{to:"/events",color:"light-blue"}},[r("v-list-item-content",[r("v-list-item-title",[t._v("Events")])],1)],1)],1),r("v-divider")],1),r("v-main",{staticClass:"grey lighten-3"},[r("router-view")],1)],1)},l=[],c=r("2f62"),u=r("ade3"),p=r("1da1"),h=(r("96cf"),r("bc3a")),d=r.n(h),v=r("3835");r("d81d"),r("13d5"),r("4fad"),r("d3b7"),r("25f0"),r("c1f9");(function(t){t["NONE"]="none",t["DEBUG"]="debug",t["INFORMATIONAL"]="informational",t["NOTICE"]="notice",t["WARNING"]="warning",t["ERROR"]="error",t["CRITICAL"]="critical",t["ALERT"]="alert",t["EMERGENCY"]="emergency"})(a||(a={}));var m,f=function(t){return t.map((function(t){var e=Object.entries(t.output_fields||{}).reduce((function(t,e){var r=Object(v["a"])(e,2),a=r[0],i=r[1];return t[a]=new String(i).toString(),t}),{}),r=Object.entries(e);return r.sort((function(t,e){var r=Object(v["a"])(t,2),a=r[0],i=r[1],n=Object(v["a"])(e,2),s=n[0],o=n[1];return a.length+i.length-(s.length+o.length)})),{output:t.output,priority:t.priority.toLowerCase(),rule:t.rule,time:new Date(t.time),outputFields:Object.fromEntries(r),filterTime:1e3*Math.floor(new Date(t.time).getTime()/1e3)}})).reverse()},b=function(){var t;return t={},Object(u["a"])(t,a.EMERGENCY,0),Object(u["a"])(t,a.ALERT,0),Object(u["a"])(t,a.CRITICAL,0),Object(u["a"])(t,a.ERROR,0),Object(u["a"])(t,a.WARNING,0),Object(u["a"])(t,a.NOTICE,0),Object(u["a"])(t,a.INFORMATIONAL,0),Object(u["a"])(t,a.DEBUG,0),Object(u["a"])(t,a.NONE,0),t},y=(i={},Object(u["a"])(i,a.NONE,"grey"),Object(u["a"])(i,a.DEBUG,"light-blue lighten-1"),Object(u["a"])(i,a.INFORMATIONAL,"light-blue"),Object(u["a"])(i,a.NOTICE,"primary"),Object(u["a"])(i,a.WARNING,"warning"),Object(u["a"])(i,a.ERROR,"error"),Object(u["a"])(i,a.CRITICAL,"red darken-1"),Object(u["a"])(i,a.ALERT,"red darken-2"),Object(u["a"])(i,a.EMERGENCY,"red darken-3"),i),O=(n={total:"#757575"},Object(u["a"])(n,a.EMERGENCY,"#C62828"),Object(u["a"])(n,a.ALERT,"#D32F2F"),Object(u["a"])(n,a.CRITICAL,"#E53935"),Object(u["a"])(n,a.ERROR,"#FF5252"),Object(u["a"])(n,a.WARNING,"#FB8C00"),Object(u["a"])(n,a.NOTICE,"#1976D2"),Object(u["a"])(n,a.INFORMATIONAL,"#03A9F4"),Object(u["a"])(n,a.DEBUG,"#29B6F6"),Object(u["a"])(n,a.NONE,"#555"),n),g=function(t){return y[t]},C=!0,x=d.a.create({baseURL:C?"//".concat(window.location.host,"/"):"http://localhost:2802"}),j={event:function(){return Object(p["a"])(regeneratorRuntime.mark((function t(){var e,r;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.next=2,x.get("events");case 2:return e=t.sent,r=e.data,t.abrupt("return",{events:f(r.events||[]),retention:r.retention,stats:r.stats,outputs:r.outputs});case 5:case"end":return t.stop()}}),t)})))()}},_=r("f7fe"),w=r.n(_);s["a"].use(c["a"]);var E="FETCH_EVENTS",T="RECEIVE_EVENTS",R="SET_STATE",S=new c["a"].Store({state:{socket:{isConnected:!1},events:[],outputs:[],retention:0},mutations:Object(u["a"])({},R,(function(t,e){t.events=e.events||[],t.outputs=e.outputs||[],t.retention=e.retention})),actions:(m={},Object(u["a"])(m,E,(function(t){var e=t.commit;j.event().then((function(t){return e(R,t)}))})),Object(u["a"])(m,T,w()((function(t,e){var r,i=t.commit,n={events:f(e.events||[]),outputs:e.outputs,retention:e.retention};e.stats&&(n.stats=(r={},Object(u["a"])(r,a.EMERGENCY,e.stats[a.EMERGENCY]),Object(u["a"])(r,a.ALERT,e.stats[a.ALERT]),Object(u["a"])(r,a.CRITICAL,e.stats[a.CRITICAL]),Object(u["a"])(r,a.ERROR,e.stats[a.ERROR]),Object(u["a"])(r,a.WARNING,e.stats[a.WARNING]),Object(u["a"])(r,a.NOTICE,e.stats[a.NOTICE]),Object(u["a"])(r,a.INFORMATIONAL,e.stats[a.INFORMATIONAL]),Object(u["a"])(r,a.DEBUG,e.stats[a.DEBUG]),Object(u["a"])(r,a.NONE,e.stats[a.NONE]),r));i(R,n)}),1e3)),m)}),A=s["a"].extend({name:"App",computed:Object(c["b"])(["retention","outputs"]),created:function(){this.$store.dispatch(E)}}),k=A,N=r("2877"),V=r("6544"),$=r.n(V),I=r("7496"),L=r("cc20"),M=r("ce7e"),P=r("adda"),F=r("8860"),q=r("da13"),D=r("5d23"),G=r("f6c4"),z=r("f774"),B=r("8dd9"),H=Object(N["a"])(k,o,l,!1,null,null,null),W=H.exports;$()(H,{VApp:I["a"],VChip:L["a"],VDivider:M["a"],VImg:P["a"],VList:F["a"],VListItem:q["a"],VListItemContent:D["a"],VListItemTitle:D["b"],VMain:G["a"],VNavigationDrawer:z["a"],VSheet:B["a"]});var U=r("8c4f"),Y=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-container",{staticClass:"px-3 py-2",attrs:{fluid:""}},[r("v-row",{staticClass:"pb-0"},[r("v-col",{attrs:{cols:"12"}},[r("v-card",{attrs:{elevation:"2"}},[r("v-toolbar",{attrs:{flat:""}},[r("v-spacer"),r("rule-autocomplete",{staticClass:"ml-3",staticStyle:{"max-width":"500px"},attrs:{events:t.events},model:{value:t.rules,callback:function(e){t.rules=e},expression:"rules"}}),r("priority-autocomplete",{staticClass:"ml-3",staticStyle:{"max-width":"350px"},model:{value:t.priorities,callback:function(e){t.priorities=e},expression:"priorities"}}),r("time-range-select",{staticClass:"ml-3",staticStyle:{"max-width":"300px"},model:{value:t.time,callback:function(e){t.time=e},expression:"time"}})],1)],1)],1)],1),r("v-row",{staticClass:"pb-0 mt-0"},[r("v-col",{attrs:{cols:"4"}},[r("priority-chart",{attrs:{events:t.filtered}})],1),r("v-col",{attrs:{cols:"8"}},[r("rule-chart",{attrs:{events:t.filtered}})],1)],1),r("v-row",{staticClass:"pb-0 mt-0"},[r("v-col",{attrs:{cols:"12"}},[r("rule-time-chart",{key:t.time,attrs:{events:t.filtered,time:t.time}})],1)],1),r("v-row",{staticClass:"pb-0 mt-0"},[r("v-col",{attrs:{cols:"12"}},[r("priority-time-chart",{key:t.time,attrs:{events:t.filtered,time:t.time}})],1)],1)],1)},J=[],K=r("5530"),Q=(r("b0c0"),r("4de4"),function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-autocomplete",t._b({attrs:{dense:"",multiple:"",items:t.options,outlined:"","hide-details":"",clearable:"",name:"Priorities",label:"Priorities",value:t.value},on:{input:t.input},scopedSlots:t._u([{key:"selection",fn:function(e){var a=e.item,i=e.index;return[i<=1?r("v-chip",{attrs:{small:"",label:"",outlined:""}},[r("span",[t._v(t._s(a))])]):t._e(),2===i?r("span",{staticClass:"grey--text caption"},[t._v(" (+"+t._s(t.selected.length-2)+" others) ")]):t._e()]}}])},"v-autocomplete",t.$attrs,!1))}),X=[],Z=(r("b64b"),w()((function(t){t()}),800)),tt=s["a"].extend({props:{value:{type:Array,default:function(){return[]}}},data:function(){return{options:Object.keys(b()),selected:this.value}},methods:{input:function(t){var e=this;this.selected=t,Z((function(){e.$emit("input",t)}))}}}),et=tt,rt=r("c6a6"),at=Object(N["a"])(et,Q,X,!1,null,null,null),it=at.exports;$()(at,{VAutocomplete:rt["a"],VChip:L["a"]});var nt=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-card",{staticStyle:{height:"100%"}},[r("v-card-title",[t._v(" Priorities ")]),r("v-card-text",[r("wait",[r("apexchart",{attrs:{type:"donut",height:"350",options:t.pie.chartOptions,series:t.pie.series}})],1)],1)],1)},st=[],ot=(r("a9e3"),r("159b"),r("07ac"),function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("div",[t.ready?t._t("default"):t._e()],2)}),lt=[],ct=s["a"].extend({data:function(){return{ready:!1}},created:function(){var t=this;setTimeout((function(){t.ready=!0}),100)}}),ut=ct,pt=Object(N["a"])(ut,ot,lt,!1,null,null,null),ht=pt.exports,dt=s["a"].extend({components:{Wait:ht},name:"PieChart",props:{events:{type:Array,required:!0},time:{type:Number}},computed:{stats:function(){var t=b();return this.events.forEach((function(e){t[e.priority]+=1})),t},pie:function(){return{series:Object.values(this.stats),chartOptions:{chart:{height:350,type:"donut"},dataLabels:{dropShadow:{enabled:!0,top:0,left:0,right:0,bottom:0,blur:1,color:"#000",opacity:1}},plotOptions:{pie:{donut:{labels:{show:!0,total:{showAlways:!0,show:!0}}}}},labels:Object.keys(this.stats),colors:Object.keys(this.stats).map((function(t){return O[t]}))}}}}}),vt=dt,mt=r("b0af"),ft=r("99d9"),bt=Object(N["a"])(vt,nt,st,!1,null,null,null),yt=bt.exports;$()(bt,{VCard:mt["a"],VCardText:ft["a"],VCardTitle:ft["b"]});var Ot=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-card",{staticStyle:{height:"100%"}},[r("v-card-title",[t._v(" Priority Timeline ")]),r("v-card-text",[r("wait",[r("apexchart",{attrs:{type:"bar",height:t.height,options:t.chart,series:t.series}})],1)],1)],1)},gt=[],Ct=r("2909"),xt=(r("a630"),r("3ca3"),r("ddb0"),r("fb6a"),r("5319"),r("ac1f"),r("1276"),function(t,e){var r=(new Date).getTime()-t;return r-r%e}),jt=function(t){var e=arguments.length>1&&void 0!==arguments[1]?arguments[1]:12,r=t||3e5,a=r/e;return{range:r,steps:e,start:xt(r,a),stepSize:a}},_t=s["a"].extend({components:{Wait:ht},name:"PriorityTimeChart",props:{events:{type:Array,required:!0},time:{type:Number}},data:function(){var t=jt(this.time,15),e=t.range,r=t.steps,a=t.start,i=t.stepSize;return{range:e,steps:r,start:a,stepSize:i,height:400}},watch:{events:function(){this.start=xt(this.range,this.stepSize)}},computed:{categories:function(){var t=this;return Array.from(Array(this.steps).keys()).reduce((function(e,r){var a=t.start+(r+1)*t.stepSize,i=new Date(a).toISOString().replace("T"," ").slice(0,-5);return t.range>=36e5&&(i="".concat(i.slice(0,-2),"00")),e.push({time:a,label:i}),e}),[])},series:function(){var t=this,e=Object.keys(b()).reduce((function(e,r){return e[r]={name:r,data:Object(Ct["a"])(Array(t.steps).keys()).map((function(){return 0}))},e}),{}),r=this.events.reduce((function(e,r){if(r.filterTime<t.start)return e;for(var a in t.categories)if(t.categories[+a+1]&&r.filterTime<t.categories[+a+1].time)return e[r.priority].data[a]+=1,e;return e[r.priority].data[t.steps-1]+=1,e}),e);return Object.values(r).reverse()},chart:function(){return{chart:{type:"bar",height:this.height,stacked:!0,toolbar:{show:!1}},colors:this.series.map((function(t){return O[t.name]})),plotOptions:{bar:{horizontal:!1}},xaxis:{categories:this.categories.map((function(t){var e=t.label;return e.split(" ")}))},legend:{position:"right",offsetY:40},fill:{opacity:1}}}}}),wt=_t,Et=Object(N["a"])(wt,Ot,gt,!1,null,null,null),Tt=Et.exports;$()(Et,{VCard:mt["a"],VCardText:ft["a"],VCardTitle:ft["b"]});var Rt=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-autocomplete",t._b({attrs:{dense:"",multiple:"",items:t.options,outlined:"","hide-details":"",name:"Rules",label:"Rules",clearable:"",value:t.value},on:{input:t.input},scopedSlots:t._u([{key:"selection",fn:function(e){var a=e.item,i=e.index;return[0===i?r("v-chip",{attrs:{small:"",label:"",outlined:""}},[r("span",[t._v(t._s(a))])]):t._e(),1===i?r("span",{staticClass:"grey--text caption"},[t._v(" (+"+t._s(t.value.length-1)+" others) ")]):t._e()]}}])},"v-autocomplete",t.$attrs,!1))},St=[],At=w()((function(t){t()}),800),kt=s["a"].extend({props:{events:{type:Array,default:function(){return[]}},value:{type:Array,default:function(){return[]}}},data:function(){return{selected:[]}},computed:{options:function(){return Object.keys(this.events.reduce((function(t,e){return t[e.rule]=null,t}),{})).sort()}},methods:{input:function(t){var e=this;this.selected=t,At((function(){e.$emit("input",t)}))}}}),Nt=kt,Vt=Object(N["a"])(Nt,Rt,St,!1,null,null,null),$t=Vt.exports;$()(Vt,{VAutocomplete:rt["a"],VChip:L["a"]});var It=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-card",{staticStyle:{height:"100%"}},[r("v-card-title",[t._v(" Rules ")]),r("v-card-text",[r("apexchart",{attrs:{type:"bar",height:"350",options:t.chart.chartOptions,series:t.chart.series}})],1)],1)},Lt=[],Mt=(r("a15b"),function(t){return t.split(" ").reduce((function(t,e,r){var a=Math.floor(r/3);return t[a]||(t[a]=[]),t[a].push(e),t}),[]).map((function(t){return t.join(" ")}))}),Pt=s["a"].extend({name:"RuleChart",props:{events:{type:Array,required:!0},time:{type:Number}},computed:{rules:function(){var t=this.events.reduce((function(t,e){return!1===Object.prototype.hasOwnProperty.call(t,e.rule)&&(t[e.rule]=0),t[e.rule]+=1,t}),{}),e=Object.keys(t).sort().reduce((function(e,r){return e[r]=t[r],e}),{});return e},chart:function(){return{series:[{name:"Counter",data:Object.values(this.rules)}],chartOptions:{chart:{height:350,type:"bar",toolbar:{show:!1},zoom:{enabled:!1}},dataLabels:{enabled:!1},xaxis:{labels:{rotate:-45,rotateAlways:!0,hideOverlappingLabels:!1,maxHeight:200,style:{colors:[],fontSize:"12px"}},categories:Object.keys(this.rules).map((function(t){return Mt(t)})),tickPlacement:"on"}}}}}}),Ft=Pt,qt=Object(N["a"])(Ft,It,Lt,!1,null,null,null),Dt=qt.exports;$()(qt,{VCard:mt["a"],VCardText:ft["a"],VCardTitle:ft["b"]});var Gt=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-card",{staticStyle:{height:"100%"}},[r("v-card-title",[t._v(" Rule Timeline ")]),r("v-card-text",[r("wait",[r("apexchart",{attrs:{type:"bar",height:t.height,options:t.chart,series:t.series}})],1)],1)],1)},zt=[],Bt=s["a"].extend({components:{Wait:ht},name:"RuleTimeChart",props:{events:{type:Array,required:!0},time:{type:Number}},data:function(){var t=jt(this.time,15),e=t.range,r=t.steps,a=t.start,i=t.stepSize;return{range:e,steps:r,start:a,stepSize:i,height:400}},watch:{events:function(){this.start=xt(this.range,this.stepSize)}},computed:{categories:function(){var t=this;return Array.from(Array(this.steps).keys()).reduce((function(e,r){var a=t.start+(r+1)*t.stepSize,i=new Date(a).toISOString().replace("T"," ").slice(0,-5);return t.range>=36e5&&(i="".concat(i.slice(0,-2),"00")),e.push({time:a,label:i}),e}),[])},series:function(){var t=this,e=this.events.reduce((function(e,r){if(!1===Object.prototype.hasOwnProperty.call(e,r.rule)&&(e[r.rule]={name:r.rule,data:Object(Ct["a"])(Array(t.steps).keys()).map((function(){return 0}))}),r.filterTime<t.start)return e;for(var a in t.categories)if(t.categories[+a+1]&&r.filterTime<t.categories[+a+1].time)return e[r.rule].data[a]+=1,e;return e[r.rule].data[t.steps-1]+=1,e}),{}),r=Object.keys(e).sort().reduce((function(t,r){return t[r]=e[r],t}),{});return Object.values(r).filter((function(t){var e=t.data.reduce((function(t,e){return t+e}),0);return e>0}))},chart:function(){return{chart:{type:"bar",height:this.height,stacked:!0,toolbar:{show:!1},zoom:{enabled:!1}},plotOptions:{bar:{horizontal:!1}},xaxis:{categories:this.categories.map((function(t){var e=t.label;return e.split(" ")}))},legend:{position:"right",offsetY:40},fill:{opacity:1}}}}}),Ht=Bt,Wt=Object(N["a"])(Ht,Gt,zt,!1,null,null,null),Ut=Wt.exports;$()(Wt,{VCard:mt["a"],VCardText:ft["a"],VCardTitle:ft["b"]});var Yt=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-select",t._g(t._b({attrs:{outlined:"","hide-details":"",dense:"",label:"Timerange","item-value":"value","item-text":"label",items:t.options,clearable:""}},"v-select",t.$attrs,!1),t.$listeners))},Jt=[],Kt=s["a"].extend({inheritAttrs:!1,data:function(){return{options:[{label:"last 15 minutes",value:9e5},{label:"last 30 minutes",value:18e5},{label:"last hour",value:36e5},{label:"last 6 hours",value:216e5},{label:"last 12 hours",value:432e5},{label:"last 24 hours",value:864e5},{label:"last 2 days",value:1728e5},{label:"last 7 days",value:6048e5}]}}}),Qt=Kt,Xt=r("b974"),Zt=Object(N["a"])(Qt,Yt,Jt,!1,null,null,null),te=Zt.exports;$()(Zt,{VSelect:Xt["a"]});var ee=s["a"].extend({components:{TimeRangeSelect:te,PriorityChart:yt,RuleChart:Dt,RuleTimeChart:Ut,PriorityTimeChart:Tt,PriorityAutocomplete:it,RuleAutocomplete:$t},name:"Dashboard",computed:Object(K["a"])(Object(K["a"])({},Object(c["b"])(["events"])),{},{time:{get:function(){var t=this.$route.query.time;return t?+t:null},set:function(t){t?this.$router.push({name:this.$route.name||"",query:Object(K["a"])(Object(K["a"])({},this.$route.query),{},{time:t.toString()})}):this.$router.push({name:this.$route.name||"",query:Object(K["a"])(Object(K["a"])({},this.$route.query),{},{time:void 0})})}},priorities:{get:function(){var t=this.$route.query.priority;return t?Array.isArray(t)?t:[t]:[]},set:function(t){t?this.$router.push({name:this.$route.name||"",query:Object(K["a"])(Object(K["a"])({},this.$route.query),{},{priority:t})}):this.$router.push({name:this.$route.name||"",query:Object(K["a"])(Object(K["a"])({},this.$route.query),{},{priority:void 0})})}},rules:{get:function(){var t=this.$route.query.rules;return t?Array.isArray(t)?t:[t]:[]},set:function(t){t?this.$router.push({name:this.$route.name||"",query:Object(K["a"])(Object(K["a"])({},this.$route.query),{},{rules:t})}):this.$router.push({name:this.$route.name||"",query:Object(K["a"])(Object(K["a"])({},this.$route.query),{},{rules:void 0})})}},filtered:function(){var t=this,e=this.time,r=null;return"number"===typeof e&&e>0&&(r=(new Date).getTime()-e),0!==this.rules.length||0!==this.priorities.length||r?this.events.filter((function(e){return(!r||e.filterTime>=r)&&(!t.priorities.length||t.priorities.includes(e.priority))&&(!t.rules.length||t.rules.includes(e.rule))})):this.events}})}),re=ee,ae=r("62ad"),ie=r("a523"),ne=r("0fd9"),se=r("2fa4"),oe=r("71d9"),le=Object(N["a"])(re,Y,J,!1,null,null,null),ce=le.exports;$()(le,{VCard:mt["a"],VCol:ae["a"],VContainer:ie["a"],VRow:ne["a"],VSpacer:se["a"],VToolbar:oe["a"]});var ue=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-container",{staticClass:"px-3 py-2",attrs:{fluid:""}},[r("v-row",{staticClass:"pb-0"},[r("v-col",{attrs:{cols:"12"}},[r("v-card",{attrs:{elevation:"2"}},[r("v-card-text",t._l(t.stats,(function(e,a){return r("priority-counter-chip",{key:a,attrs:{dark:"",priority:a,count:e},on:{click:function(e){t.search=a}}},[t._v(" "+t._s(a)+" "),r("div",{staticClass:"white black--text rounded ml-2 px-2 py-2 text-body-2"},[t._v(t._s(e))])])})),1)],1)],1)],1),r("v-row",{staticClass:"pt-0 mt-0"},[r("v-col",{attrs:{cols:"12"}},[r("v-card",{attrs:{elevation:"2"}},[r("v-toolbar",{attrs:{flat:""}},[r("v-spacer"),r("time-range-select",{staticClass:"ml-3",staticStyle:{"max-width":"300px"},model:{value:t.time,callback:function(e){t.time=e},expression:"time"}}),r("v-text-field",{staticClass:"ml-3",staticStyle:{"max-width":"300px"},attrs:{outlined:"","hide-details":"",dense:"",label:"Search",clearable:""},model:{value:t.search,callback:function(e){t.search=e},expression:"search"}})],1),r("v-divider"),r("event-table",{attrs:{events:t.filtered,search:t.search},on:{"update:search":function(e){t.search=e}}})],1)],1)],1)],1)},pe=[],he=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-data-table",{attrs:{items:t.events,headers:t.headers,"items-per-page":-1,"hide-default-footer":"",search:t.search,"custom-filter":t.defaultFilter,"disable-sort":""},scopedSlots:t._u([{key:"item",fn:function(e){var a=e.item;return[r("tr",[r("td",[r("span",{domProps:{innerHTML:t._s(a.time.toISOString().split("T").join(" ").slice(0,-5))}})]),r("td",[r("v-chip",{staticClass:"mb-3 white--text",attrs:{label:"",color:t.color(a.priority)},on:{click:function(e){return t.$emit("update:search",a.priority)}}},[r("span",{domProps:{innerHTML:t._s(t.highlightMatches(a.priority))}})]),r("br"),r("v-chip",{staticClass:"black--text",attrs:{label:"",color:"light-blue lighten-5"},domProps:{innerHTML:t._s(t.highlightMatches(a.rule))},on:{click:function(e){return t.$emit("update:search",a.rule)}}}),r("br")],1),r("td",{staticClass:"py-3",staticStyle:{"font-size":"0.8rem"}},[r("div",{domProps:{innerHTML:t._s(t.highlightMatches(a.output))}}),r("div",{staticClass:"mt-3"},[t._l(a.outputFields,(function(e,a){return[e.length<=100?r("output-field-chip",{key:a,attrs:{dark:"",value:t.highlightMatches(e),label:a},on:{click:function(r){return t.$emit("update:search",e)}}}):r("output-field-card",{key:a,attrs:{value:t.highlightMatches(e),label:a},on:{click:function(r){return t.$emit("update:search",e)}}})]}))],2)])])]}}])})},de=[],ve=(r("841c"),r("4d63"),function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-card",t._g(t._b({staticClass:"primary mt-2",staticStyle:{border:"1px solid"},attrs:{flat:""}},"v-card",t.$attrs,!1),t.$listeners),[r("v-system-bar",{attrs:{color:"primary white--text"}},[t._v(t._s(t.label))]),r("v-card-text",{staticClass:"pa-2 white black--text",staticStyle:{"font-size":"0.8rem","line-height":"1.1rem"}},[r("span",{domProps:{innerHTML:t._s(t.value)}})])],1)}),me=[],fe=s["a"].extend({name:"OutputFieldCard",inheritAttrs:!1,props:{label:{type:String,required:!0},value:{type:String,required:!0}}}),be=fe,ye=r("afd9"),Oe=Object(N["a"])(be,ve,me,!1,null,null,null),ge=Oe.exports;$()(Oe,{VCard:mt["a"],VCardText:ft["a"],VSystemBar:ye["a"]});var Ce=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-chip",t._g({staticClass:"my-1 mr-2 pr-0 py-0",staticStyle:{border:"1px solid #000"},attrs:{dark:"",label:"",small:"",color:"primary"}},t.$listeners),[t._v(" "+t._s(t.label)+" "),r("div",{staticClass:"white black--text rounded ml-2 px-2 py-2 text-body-2",staticStyle:{"font-size":"0.7rem!important","line-height":"0.9rem!important"}},[r("span",{domProps:{innerHTML:t._s(t.value)}})])])},xe=[],je=s["a"].extend({name:"OutputFieldChip",props:{label:{type:String,required:!0},value:{type:String,required:!0}}}),_e=je,we=Object(N["a"])(_e,Ce,xe,!1,null,null,null),Ee=we.exports;$()(we,{VChip:L["a"]});var Te=s["a"].extend({components:{OutputFieldChip:Ee,OutputFieldCard:ge},name:"EventTable",props:{events:{type:Array,default:function(){return[]}},search:{type:String}},data:function(){return{headers:[{text:"Time",value:"time"},{text:"Rule / Priority",value:"rule"},{text:"Output",value:"output",width:"60%"}]}},methods:{color:function(t){return g(t)},highlightMatches:function(t){if(!this.search)return t;var e=t.toLowerCase().includes(this.search.toLowerCase());return e?t.replace(new RegExp(this.search,"ig"),(function(t){return"<strong>".concat(t,"</strong>")})):t},defaultFilter:function(t,e,r){e=(e||"").toLocaleLowerCase();var a=Object.values(r.outputFields).some((function(t){return t.toLowerCase().includes(e)}));if(a)return!0;var i=null!=r.priority&&null!=e&&-1!==r.priority.toLocaleLowerCase().indexOf(e);return!!i||null!=t&&null!=e&&"boolean"!==typeof t&&-1!==t.toString().toLocaleLowerCase().indexOf(e)}}}),Re=Te,Se=(r("fedc"),r("8fea")),Ae=Object(N["a"])(Re,he,de,!1,null,"2402e04a",null),ke=Ae.exports;$()(Ae,{VChip:L["a"],VDataTable:Se["a"]});var Ne=function(){var t=this,e=t.$createElement,r=t._self._c||e;return r("v-chip",t._g({staticClass:"my-0 mr-2 pr-0 py-0",staticStyle:{border:"1px solid #000"},attrs:{dark:"",color:t.color(t.priority),label:""}},t.$listeners),[t._v(" "+t._s(t.priority)+" "),r("div",{staticClass:"white black--text rounded ml-2 px-2 py-2 text-body-2"},[t._v(t._s(t.count))])])},Ve=[],$e=s["a"].extend({name:"PriorityCountChip",props:{priority:{type:String,required:!0},count:{type:Number,required:!0}},methods:{color:function(t){return g(t)}}}),Ie=$e,Le=Object(N["a"])(Ie,Ne,Ve,!1,null,null,null),Me=Le.exports;$()(Le,{VChip:L["a"]});var Pe=s["a"].extend({components:{EventTable:ke,PriorityCounterChip:Me,TimeRangeSelect:te},name:"Dashboard",data:function(){return{search:"",time:null}},computed:Object(K["a"])(Object(K["a"])({},Object(c["b"])(["events"])),{},{filtered:function(){var t=this.time;if("number"===typeof t&&t>0){var e=(new Date).getTime()-t;return this.events.filter((function(t){return t.filterTime>=e}))}return this.events},stats:function(){var t=Object(K["a"])({total:0},b());return this.events.forEach((function(e){t[e.priority]+=1,t.total+=1})),t}}),methods:{color:function(t){return g(t)}}}),Fe=Pe,qe=r("8654"),De=Object(N["a"])(Fe,ue,pe,!1,null,null,null),Ge=De.exports;$()(De,{VCard:mt["a"],VCardText:ft["a"],VCol:ae["a"],VContainer:ie["a"],VDivider:M["a"],VRow:ne["a"],VSpacer:se["a"],VTextField:qe["a"],VToolbar:oe["a"]}),s["a"].use(U["a"]);var ze=[{path:"/",name:"Dashboard",component:ce},{path:"/events",name:"Events",component:Ge}],Be=new U["a"]({mode:"hash",base:"/ui/",routes:ze}),He=Be,We=r("f309");s["a"].use(We["a"]);var Ue=new We["a"]({}),Ye=r("b408"),Je=r.n(Ye),Ke=r("1321"),Qe=r.n(Ke);s["a"].use(Qe.a),s["a"].component("apexchart",Qe.a);var Xe="ws://".concat(window.location.host,"/ws");s["a"].use(Je.a,Xe,{format:"json",reconnection:!0,store:S,passToStoreHandler:function(t,e){t.includes("onmessage")&&this.store.dispatch(T,JSON.parse(e.data))}}),s["a"].config.productionTip=!1,new s["a"]({router:He,store:S,vuetify:Ue,render:function(t){return t(W)}}).$mount("#app")},fedc:function(t,e,r){"use strict";r("4080")}});
//# sourceMappingURL=app.23760645.js.map