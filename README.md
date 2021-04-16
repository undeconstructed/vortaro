# Vortaro

Provo krei novan Esperantan vortaron el la Akademia Vortaro (AV), kaj indiki kiel uzebla estas ĉiu vorto.

## Fontoj

f = fundamento
a = aldono
n = neoficiala

## Specoj

### e = esperanta

kernaj partoj de la lingvo.

### a = abstrakta

se tempo subite haltus, abstraktaj aferoj ne plu ekzistos.
tuj uzeblaj kiel verboj, adjektivoj, adverboj - la substantivo kutima estas la ago aŭ la rezulto.
ĝenerale estas la plej potenca, ĉar oni povas libere uzi ul', aĵ', il', ist' ktp..

### k = konkreta

se tempo subite haltus, konkretaj aferoj plu ekzistos.
la substantivo estas la aĵo mem.
la aliaj formoj signifias "esti kiel la aĵo", aŭ "fari per la aĵo".

### n = nomoj

nomo por unu el kelkaj aferoj da sama speco.

## Niveloj

0 = ĉiuj komprenos, ekz. "sub"
1 = ofta vorto, ĝenerale komprenebla, ekz. "striki"
2 = faka vorto, kiu eble ne estas ĝenerale konata, ekz. "sturgo"
3 = neniu komprenus sen speciala intereso, ekz. "stofo", "stupo"

## Bone?

Ĉu la vorto havas gravajn problemojn, ekz. ĝi ne havas veran radikon, kiel "teism'".

## Komandoj

Ĉiuj fakoj:

`cat radikoj.json | jq -r '.[] .fakoj[]' | sort | uniq`

Priskribataj fakoj:

`cat fakoj.json | jq -r '.[] .nomo' | sort`

Vortoj laŭ fako:

`cat radikoj.json | jq -r '.[] | select(.fakoj | contains(["ĝenerala"])) | .radiko' | sed 's/\///g'`

Vortoj uzataj en ekzemploj:

`cat ekzemploj.json | jq -r '.[]' | tr '[A-ZĈĜĤĴŜŬ]' '[a-zĉĝĥĵŝŭ]' | sed 's/[^a-zĉĝĥĵŝŭ]/\n/g' | sort | uniq`

Vortoj sen ekzemploj:

`comm -13 <(cat ekzemploj.json | jq -r '.[]' | tr '[A-ZĈĜĤĴŜŬ]' '[a-zĉĝĥĵŝŭ]' | sed 's/[^a-zĉĝĥĵŝŭ]/\n/g' | sort | uniq) <(cat radikoj.json | jq -r '.[] | select(.fakoj | contains(["ĝenerala"])) | .radiko' | sed 's/\///g' | sort)`
