package blog

const Host = "https://www.notion.so/"

const HugoBlogHeaderTpl = `---
showonlyimage: true
title:      "{{.Title}}"
subtitle:   "{{.SubTitle}}"
excerpt: "{{.Excerpt}}"
description: "{{.Description}}"
date:       "{{.Date}}"
author: "foreversmart"
image: "{{.Image}}"
published: true 
tags: [{{range .Tags}}{{.}}, {{end}}]
categories: [{{range .Category}}{{.}}, {{end}}]
URL: "/{{.UrlPath}}/{{.Url}}/"
---`
