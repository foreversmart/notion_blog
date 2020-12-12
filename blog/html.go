package blog

const Tail = `</body>
</html>
`
const Mid = `
<div style="width: 100%; display: flex; flex-direction: column; align-items: center; flex-shrink: 0; flex-grow: 0;"><div style="padding-left: calc(96px + env(safe-area-inset-left)); padding-right: calc(96px + env(safe-area-inset-right)); max-width: 100%; width: 900px;"><div class="notion-page-details-controls" style="display: flex; align-items: baseline; justify-content: flex-start; flex-wrap: wrap; color: rgba(55, 53, 47, 0.4); font-family: -apple-system, BlinkMacSystemFont, &quot;Segoe UI&quot;, Helvetica, &quot;Apple Color Emoji&quot;, Arial, sans-serif, &quot;Segoe UI Emoji&quot;, &quot;Segoe UI Symbol&quot;; padding-bottom: 20px;"></div></div></div>
`
const TitleStyle = `
<div style="padding-left: calc(96px + env(safe-area-inset-left)); padding-right: calc(96px + env(safe-area-inset-right)); max-width: 100%; width: 900px;">`

const ContentStyle = ``

const Head = `
<html>
<head>
    <title></title>
</head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<style>body {
        background: #fff
    }

    body.dark {
        background: #2f3437
    }

	code[class*="language-"],
pre[class*="language-"] {
	color: black;
	background: none;
	text-shadow: 0 1px white;
	font-family: Consolas, Monaco, 'Andale Mono', 'Ubuntu Mono', monospace;
	font-size: 1em;
	text-align: left;
	white-space: pre;
	word-spacing: normal;
	word-break: normal;
	word-wrap: normal;
	line-height: 1.5;

	-moz-tab-size: 4;
	-o-tab-size: 4;
	tab-size: 4;

	-webkit-hyphens: none;
	-moz-hyphens: none;
	-ms-hyphens: none;
	hyphens: none;
}
pre[class*="language-"]::-moz-selection, pre[class*="language-"] ::-moz-selection, code[class*="language-"]::-moz-selection, code[class*="language-"] ::-moz-selection {
	text-shadow: none;
	background: #b3d4fc;
}
pre[class*="language-"]::selection, pre[class*="language-"] ::selection, code[class*="language-"]::selection, code[class*="language-"] ::selection {
	text-shadow: none;
	background: #b3d4fc;
}
@media print {
	code[class*="language-"],
	pre[class*="language-"] {
		text-shadow: none;
	}
}
/* Code blocks */
pre[class*="language-"] {
	padding: 1em;
	margin: .5em 0;
	overflow: auto;
}
:not(pre) > code[class*="language-"],
pre[class*="language-"] {
	background: #f5f2f0;
}
/* Inline code */
:not(pre) > code[class*="language-"] {
	padding: .1em;
	border-radius: .3em;
	white-space: normal;
}
.token.comment,
.token.prolog,
.token.doctype,
.token.cdata {
	color: slategray;
}
.token.punctuation {
	color: #999;
}
.token.namespace {
	opacity: .7;
}
.token.property,
.token.tag,
.token.boolean,
.token.number,
.token.constant,
.token.symbol,
.token.deleted {
	color: #905;
}
.token.selector,
.token.attr-name,
.token.string,
.token.char,
.token.builtin,
.token.inserted {
	color: #690;
}
.token.operator,
.token.entity,
.token.url,
.language-css .token.string,
.style .token.string {
	color: #9a6e3a;
	background: hsla(0, 0%, 100%, .5);
}
.token.atrule,
.token.attr-value,
.token.keyword {
	color: #07a;
}
.token.function,
.token.class-name {
	color: #DD4A68;
}
.token.regex,
.token.important,
.token.variable {
	color: #e90;
}
.token.important,
.token.bold {
	font-weight: bold;
}
.token.italic {
	font-style: italic;
}
.token.entity {
	cursor: help;
}

/**
 * prism.js Dark theme for JavaScript, CSS and HTML
 * Based on the slides of the talk “/Reg(exp){2}lained/”
 * @author Lea Verou
 */

.notion-dark-theme code[class*="language-"],
.notion-dark-theme pre[class*="language-"] {
	color: white;
	background: none;
	text-shadow: 0 -.1em .2em black;
	font-family: Consolas, Monaco, 'Andale Mono', 'Ubuntu Mono', monospace;
	font-size: 1em;
	text-align: left;
	white-space: pre;
	word-spacing: normal;
	word-break: normal;
	word-wrap: normal;
	line-height: 1.5;

	-moz-tab-size: 4;
	-o-tab-size: 4;
	tab-size: 4;

	-webkit-hyphens: none;
	-moz-hyphens: none;
	-ms-hyphens: none;
	hyphens: none;
}

    .initial-loading-spinner {
        -webkit-animation: rotate 1s linear infinite;
        animation: rotate 1s linear infinite;
        -webkit-transform-origin: center center;
        transform-origin: center center;
        width: 1em;
        height: 1em;
        opacity: .5;
        display: block;
        pointer-events: none
    }

    @-webkit-keyframes rotate {
        0% {
            -webkit-transform: rotate(0) translateZ(0);
            transform: rotate(0) translateZ(0)
        }
        100% {
            -webkit-transform: rotate(360deg) translateZ(0);
            transform: rotate(360deg) translateZ(0)
        }
    }

    * {
        box-sizing: border-box;
    }

    @keyframes rotate {
        0% {
            -webkit-transform: rotate(0) translateZ(0);
            transform: rotate(0) translateZ(0)
        }
        100% {
            -webkit-transform: rotate(360deg) translateZ(0);
            transform: rotate(360deg) translateZ(0)
        }
    }</style>
<body>
`
