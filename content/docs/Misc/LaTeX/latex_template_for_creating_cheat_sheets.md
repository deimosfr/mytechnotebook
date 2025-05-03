---
weight: 999
url: "/Template_pour_créer_des_Cheat_Sheet_en_LaTeX/"
title: "LaTeX Template for Creating Cheat Sheets"
description: "A guide for creating cheat sheets in LaTeX with a reusable template and formatting tips."
categories: ["Linux"]
date: "2012-02-14T13:08:00+02:00"
lastmod: "2012-02-14T13:08:00+02:00"
tags: ["LaTeX", "Development", "Resume"]
toc: true
---

![LaTeX](/images/latex_logo.avif)

## Introduction

LaTeX is a language and document composition system created by Leslie Lamport in 1983. More precisely, it's a collection of macro-commands designed to facilitate the use of Donald Knuth's "text processor" TeX. Since 1993, it has been maintained by the LATEX3 Project team. The first widely used version, called LaTeX2.09, came out in 1984. A major revision, called LaTeX2ε, was released in 1991.

Having cheat sheets is often very useful when learning a new language or working with new software. I wanted to create my own cheat sheets, but finding a template wasn't easy. So I decided to share mine, which I built by drawing inspiration from others I found on the internet.

## Cheat Sheet Template

```latex {linenos=table}
% -----------------------------------------------------------------------
% Cheat Sheet Template
%
% Usage in Document content :
% \section : create a new section in column
%
% \columnbreak : break the column to force following content to jump to
%                the next column
% \cm{command}{description} : set a dotted line between the command and
%                             the description
% -----------------------------------------------------------------------

% -----------------------------------------------------------------------
% Document settings
% -----------------------------------------------------------------------

\documentclass[10pt,landscape]{article}
\usepackage{multicol}
\usepackage{calc}
\usepackage{ifthen}
\usepackage[landscape]{geometry}
\usepackage{amsmath,amsthm,amsfonts,amssymb}
\usepackage{color,graphicx,overpic}
\usepackage{hyperref}

% PDF informations
\pdfinfo{
  /Title (cheat_sheet_template.pdf)
  /Creator (TeX)
  /Producer (/pdfTeX 1.40.0)
  /Author (Pierre Mavro)
  /Subject (Example)
  /Keywords (/pdflatex, latex,pdftex,tex)}

% This sets page margins to .5 inch if using letter paper, and to 1cm
% if using A4 paper. (This probably isn't strictly necessary.)
% If using another size paper, use default 1cm margins.
\ifthenelse{\lengthtest { \paperwidth = 11in}}
    { \geometry{top=.5in,left=.5in,right=.5in,bottom=.5in} }
    {\ifthenelse{ \lengthtest{ \paperwidth = 297mm}}
        {\geometry{top=1cm,left=1cm,right=1cm,bottom=1cm} }
        {\geometry{top=1cm,left=1cm,right=1cm,bottom=1cm} }
    }

% Turn off header and footer
\pagestyle{empty}

% Redefine section commands to use less space
\makeatletter
\renewcommand{\section}{\@startsection{section}{1}{0mm}%
                                {-1ex plus -.5ex minus -.2ex}%
                                {0.5ex plus .2ex}%x
                                {\normalfont\large\bfseries}}
\renewcommand{\subsection}{\@startsection{subsection}{2}{0mm}%
                                {-1explus -.5ex minus -.2ex}%
                                {0.5ex plus .2ex}%
                                {\normalfont\normalsize\bfseries}}
\renewcommand{\subsubsection}{\@startsection{subsubsection}{3}{0mm}%
                                {-1ex plus -.5ex minus -.2ex}%
                                {1ex plus .2ex}%
                                {\normalfont\small\bfseries}}
\makeatother

% Define BibTeX command
\def\BibTeX{{\rm B\kern-.05em{\sc i\kern-.025em b}\kern-.08em
    T\kern-.1667em\lower.7ex\hbox{E}\kern-.125emX}}

% Don't print section numbers
\setcounter{secnumdepth}{0}

% Set vertical view instead of horizontal (set to 0 to let it choose)
\setcounter{unbalance}{45}

\setlength{\parindent}{0pt}
\setlength{\parskip}{0pt plus 0.5ex}

%My Environments
\newtheorem{example}[section]{Example}

% Dot lines between command and description
\def\cm#1#2{{\tt#1}\dotfill#2\par}

% -----------------------------------------------------------------------
% Document start
% -----------------------------------------------------------------------

\begin{document}
\raggedright
\footnotesize
% Set number of columns
\begin{multicols}{3}


% multicol parameters
% These lengths are set only within the two main columns
%\setlength{\columnseprule}{0.25pt}
\setlength{\premulticols}{1pt}
\setlength{\postmulticols}{1pt}
\setlength{\multicolsep}{1pt}
\setlength{\columnsep}{2pt}

\begin{center}
     \Large{\underline{Title}} \\
\end{center}

\section{Section 1}
Text
\subsection{xCode}
Subsetction text

\section{Section 2}
\cm{key}{explaination}

\section{Section 3}
Etc.

% Autor
\rule{0.3\linewidth}{0.25pt}
\section{Autor}
\href{mailto:xxx@mycompany.com}{Pierre Mavro (Deimosfr)} \\
\url{http://www.mavro.fr} \\
\url{http://www.deimos.fr}

% References
\rule{0.3\linewidth}{0.25pt}
\scriptsize
\bibliographystyle{abstract}
\bibliography{refFile}
\url{http://www.deimos.fr}
\end{multicols}
\end{document}
```

## Resources
- http://tex.stackexchange.com/questions/8827/preparing-cheat-sheets
