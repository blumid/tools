#!/bin/bash

#Usage : cat rare_subs | tech.sh

# cat /dev/stdin | httpx -silent -output httpx > origin

# get techknowledgyfile.json
webanalyze -update

webanalyze -hosts httpx -output json -silent | jq -c '{ hostname:.hostname , lang: [.matches[] | select(.app.category_names | contains(["Programming languages"]))["app_name"]]  , server: [.matches[] | select(.app.category_names | contains(["Web servers"]))["app_name"]] , cms: [.matches[] | select(.app.category_names | contains(["CMS"]))["app_name"]] }' >result

# maybe useful:
# cat result | jq '. | if .lang[0] != null then .hostname  else null end'



run_ffuf(){

    case $2 in
        php)
            wordlist=~/.cache/kiterunner/wordlists/php.txt
            ;;
        asp)
            wordlist=~/.cache/kiterunner/wordlists/httparchive_aspx_asp_cfm_svc_ashx_asmx_2022_08_28.txt
            ;;
        jsp)
            wordlist=~/.cache/kiterunner/wordlists/jsp.txt
            ;;
        all)
            wordlist=~/.cache/kiterunner/wordlists/raft-large-files.txt
            ;;
        *)
            echo "fuck!"
            exit(0)
            ;;
    esac
    wordlist=$2
    ffuf -w $wordlist -X GET -u $1/FUZZ -mc 200,405,500 -r -t 5 -ac -v -c -o "$1"_ffuf.json

    echo $1 >> f1
}

detect(){
    ## run ffuf on $1
    # ffuf -w $wordlist -X GET -u $1/FUZZ -mc 200,405,500 -r -t 5 -ac -v -c -o "$1"_ffuf.json


}

filter() {


    # case cat "result" | jq -r '. | .lang[0]' in
    #     "PHP")
    # esac

    ############## first filter #############

    ### php ###
    cat "result" | jq -r '. | select(.lang[0] == "PHP") | .hostname' | while read line; do
        echo -e "for php : \n $line"
        run_ffuf $line php
    done

    ### asp ###
    cat "result" | jq -r '. | select((.lang[0] == "ASP") or (.lang[0] == "ASPX")) | .hostname' | while read line; do
        echo -e "for asp : \n $line"
        run_ffuf $line asp
    done

    ### jsp ###
    cat "result" | jq -r '. | select(.lang[0] == "JSP") | .hostname' | while read line; do
        echo -e "for jsp : \n $line"
        run_ffuf $line jsp
    done

    # cat "result" | jq -r '. | select(.lang[0] == "PYTHON") | .hostname' | while read line; do
    #     echo -e "for python : \n $line"
    #     run_ffuf $line python
    # done

    ### node.js ###




    ############## second filter #############
    cat "origin" | anew -d "f1" | while read line; do
    # run ffuf with mixed wordlist on these
        # detect $line
        # echo $line
        run_ffuf $line all
    done

}


filter
