---
title: {{ node "title" }}
author: {{ node "author.name" }}
type: {{ getNodeType }}
date: {{ node "date" }}
url: {{ node "slug" }}
featured_image: {{ node "featuredImage.url" }}
categories: {{ node "categories" }}
tags: {{ node "tags" }}

---
{{ node "bodyRich.html" }}
