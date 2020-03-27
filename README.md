# [KATI (Koroona Abi- ja Teabe Infosüsteem)](https://kati.raulwalter.com)

*Prototüüp on püsti: https://kati.raulwalter.com*

KATI on enesediagnostikal ning andmekorjel põhinev infosüsteem mis annab vajalikku informatsiooni kriisile adekvaatseks, teadmispõhiseks reageerimiseks ning käitumisotsuste kujundamiseks ja nende operatiivseks kommunikeerimiseks kriisiolukorras.

KATI on kasutatav ka peale COVID-19 kriisi ja järgmises sarnases olukorras. Kasutatav kui Kodaniku Abi- ja Teabe Infosüsteem.

KATI lõppkasutaja- (elaniku-) poolne liides koosneb virtuaalsest enesetestist, mis suunab kasutaja läbima valikvastustest koosneva küsimustiku, et saavutada eelkõige:

- Elanikkonna tervisliku seisukorra ülevaade viiruse leviku seisukohast
- Ülevaade diagnoositud ja diagnoosimata haigusjuhtumitest
- Geograafiline paiknemine (tegelik elukoht ning külastatud paigad)
- Kasutaja lähem kontakt-ringkond läbi kontaktvõrgustiku kaardistamise

## Olulised põhimõtted
1. Kasutaja on autentitud	- Teenuse kasutamiseks peab olema kasutaja autentinud.
2. Kasutaja annab nõusoleku	- Kasutaja allkirjastab digitaalselt nõusoleku oma andmete töötlemiseks tagamaks nõnda maksimaalse vastavuse andmekaitse nõuetele.
3. Andmed on kaitstud	- Kogutud isikuandmed krüpteeritakse kasutaja autentimissertifikaadi avalikule võtmele (*täpne tehniline lahendus vajab uurimist*). Isikustamata andmeid ei krüpteerita.
4. Andmed on avatud -	Isikustamata andmed on võimalik teha kättesaadavaks Terviseametile. Mingis ulatuses ka autentimata kasutajale läbi Dashboardi.
5. Andmed kuuluvad süsteemi kasutajale -	Arendajatel/süsteemi loojatel/analüütikutel/majutajal - kellelgi ei ole ligipääsu isikustatud andmetele. Andmed kuuluvad süsteemi kasutajale. Kasutajal on võimalik igal ajal oma andmed süsteemist pöördumatult kustutada.
6. Süsteemi kasutamine on vabatahtlik.

## Kuidas KATI toimib?
KATI tuvastab kasutaja läbi eID (ID-kaart, Digi-ID, Mobiil-ID, Smart-ID, e-residendi kaart) tagamaks andmete isikustatuse ja seeläbi usaldusväärsuse. Süsteem võimaldab pidada niiöelda päevikut ka kolmanda isiku kohta (pereliige). Isiku identifikaatoriks süsteemis on isikukood. Iga kasutajal, olenemata terviseseisundist, on võimalus täita veebipõhine ankeet enda ja vajadusel oma lähikondlaste terviseseisundi kohta (sh võimalike haigussümptomite tekkimise aeg ja kirjeldus).

 Isikutuvastus | Päevik | Andmetöötlus | Tagasiside
------------- | ------------- | ------------- | ------------- 
Isikutuvastus aitab tõsta enesediagnostikal põhinevate andmete usaldusväärsust. | Kasutajal palutakse täita valikvastustega ankeet, mis võimaldab esmast, statistilistel andmetel põhinevat tervisliku seisundi analüüsi. | Andmete kogumi põhjal saab analüüsida viiruse tõenäolist levikut erinevate mudelite põhjal. Infokihti saab näha ka Eesti kaardil. | Vastustele tuginedes antakse kasutajale tagasisidet tema tervisliku seisundi kohta ning soovitused edasiseks käitumiseks.

## Täpsemalt
### Pere ja lähedaste info sisestus
Lisaks tervisliku seisundi andmetele on kasutajal võimalus ka soovi korral süsteemi lisada oma lähikondlaste nimed ja isikukoodid. Täiendavalt annab iga kasutaja enda ja oma lähikondlaste tegeliku elu -või asukoha andmed ankeedi täitmise ajahetkel ning niisamuti andmed ka mis tahes reiside kohta välisriikidesse ankeedi täitmisele eelnenud 14 päeva jooksul (Eestisse saabumise aja ja teekonna kohta). Igale kasutataja poolt loodud andmestikule tekib unikaalne identifitseerimustunnus. Kasutaja saab terviseseisundi muutumisel KATI keskkonnas andmestikku uuendada.

### KATI piirangud
Enesediagnostikal on arvestatav veamäär, kuid kuna ei ole näha, et see ajas muutuks, on tulemuseks täpsed võrdlusandmed aegribal. Teiseks saab veamäära osas kohaldada hüpoteese ning vastavalt luua erinevaid stsenaariume (konservatiivne, optimistlik).

Kui kasutaja ei ole andmeid näiteks 1 nädala jooksul uuendanud, saab ta e-kirja teel perioodilisi meeldetuletusi tervisliku seisundi ülevaatamiseks KATI keskkonnas.

### Tagasiside
Kõigile on võimaldatud ligipääs pseudonümiseeritud andmetele, kus geoinformatsioon on saadaval haldusüksuse täpsusega (linn, vald, maakond). Täiendavalt on mõeldav liidestus väliste süsteemidega võimaldamaks ligipäästu volitatud isikutele. Volitatud arstid ja avaliku sektori organisatsioonid (nt Terviseamet, haiglad, Häirekeskus, Päästeamet) saavad õigustatud vajaduse korral piiramatu ligipääsu isikustatud andmetele koos täpse geoinfoga.

KATI sisaldab muuhulgas adekvaatset informatsiooni ja ajakohaseid käitumisjuhiseid viiruse leviku tõkestamiseks. Arvestades Eestis viibivate/resideeruvate välismaalastega, on soovitatav tööriista kogutavat ja selle kaudu esitatavat teavet hoida ja anda paralleelselt nii eesti kui inglise keeles.

### Esmase triaazi võimekus
Keskkonna loomisel tuleb teha koostööd Terviseametiga kasutades nende juhiseid andmete korjamisel ja analüüsimisel. Ideaalis saab teha kogutud andmete põhjal otsuse, milliseid juhtumeid peaks edastama arstidele lähemaks uurimiseks/testimiseks ilma, et kasutaja peaks otseselt ise pöörduma oma perearsti poole. Seega leevendab keskkond riigi infotelefonide (perearstide nõuandeliin, hädaabinumber 112 ja kriisiliin 1247) kasutamisest tekkinud töökoormust.

## KATI olulisus ja mõju COVID-19 vastases võitluses

KATI toetab muuhulgas Vabariigi Valitsust selle eesmärkide saavutamisel:

- Tõkestada COVID-19 kohapealset levikut Eestis
- Tagada tervishoiusüsteemi suutlikkus COVID-19 tõrjeks ja raviks
- Vältida paanika teket ühiskonnas ja tõsta inimeste teadlikkust COVID-19 leviku tõkestamisel ja ravil.

Täpsemalt võimaldab KATI anda reaalajas Eesti Vabariigi haldusterritooriumil COVID-19 levikust süsteemsemat, interaktiivset ülevaadet. Seeläbi saab tõhustada erilolukorra operatiivset kriisikommunikatsiooni ning saavutada ühiskonnas koostöö läbi teadmispõhise teabe ja riigi suuniste.

Andmete isikustatus ning suhtlusvõrgustike ja geograafiliste asukohtade kaardistamine võimaldab meditsiinisüsteemil, sealhulgas Sotsiaalministeeriumil, Terviseametil, perearstidel, haiglatel ja samuti eriolukorra lahendamise juhtidel teha informeeritumaid otsuseid nii sihttestide planeerimiseks ja läbiviimiseks riskigrupi väliste patsientide puhul kui ka eriolukorra lahendamisel edasise nakatumisohu laienemise ärahoidmiseks.

Näiteks loob keskkonnast saadav teave eeldused otsuste kujundamiseks sotsiaalsete-, majanduslike- ja liikumispiirangute karmistamisel või leevendamisel regionaalselt vastavalt ohu tasemele. Teiseks võimaldab andmestik parendada isikukaitsevahendite kättesaadavuse vajaduse hindamist regionaalsest vaatest. Isikustatud andmete toel saaks muu hulgas määrata kohalike omavalitsuste lõikes isikuid, kes vajaks liikumispiirangutest ja/või tervislikust seisundist tulenevalt tuge näiteks esmatasandi kaupade (toit, ravimid) kojutoimetamisel.

## Olemasolevad sarnased lahendused
Eesti-sisesed veebilahendused (näiteks koroonatest.ee, koroonakaart.ee) on madala andmekvaliteediga (isikustamata, vähese detailsusega, madala kasuteguriga ohupildist ülevaate saamisel ja andmete alusel otsuste kujundamiseks) või põhinevad olemasolevatel andmetel ning ei loo uut andmenähtavust.

- Parema andmedetailsusega on Ameerika Ühendriikide Terviseameti (CDC, https://www.cdc.gov/coronavirus/2019-ncov/symptoms-testing/testing.html) analoogne lahendus, millest teiste hulgas eeskuju võtta. Plaanitav lahendus on viidatust parem, kuna tugineb isikustatud andmetele, võimaldab staatuse uuendamist/päeviku pidamist tervisenäitajate muutumisel ning sisaldab ka võimaliku nakatunu kontaktvõrgustiku, asukoha ja sagedamini külastatud geograafiliste asukohtade infot.
- COVID-19 Heat Map - https://www.evergreen-life.co.uk/covid-19-heat-map
- https://viveohealth.com/en/coronavirustest/
- https://www.cnbc.com/2020/03/25/coronavirus-singapore-to-make-contact-tracing-tech-open-source.html
- https://bluetrace.io (https://www.cnbc.com/2020/03/25/coronavirus-singapore-to-make-contact-tracing-tech-open-source.html)

## Hind
Lahenduse pakkumisel ei ole RaulWalteril ega lahenduse loomise osalistel majandushuvi ehk lahendus on Eesti Vabariigile ja kõigile kasutajatele tasuta.
