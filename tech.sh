#!/bin/bash

#Usage : cat rare_subs | tech.sh

cat /dev/stdin | httpx -silent -output httpx > origin

webanalyze -update -silent

webanalyze -hosts httpx -output json -silent | jq -c '{ hostname:.hostname , lang: [.matches[] | select(.app.category_names | contains(["Programming languages"]))["app_name"]]  , server: [.matches[] | select(.app.category_names | contains(["Web servers"]))["app_name"]] , cms: [.matches[] | select(.app.category_names | contains(["CMS"]))["app_name"]] }' >result

# maybe useful:
# cat result | jq '. | if .lang[0] != null then .hostname  else null end'

filter() {


    # case cat "result" | jq -r '. | .lang[0]' in
    #     "PHP")
    # esac

############## first filter #############

    ### php ###
    cat "result" | jq -r '. | select(.lang[0] == "PHP") | .hostname' | while read line; do
        echo -e "for php : \n $line"
        echo $line >> f1
    done

    ### asp ###
    cat "result" | jq -r '. | select(.lang[0] == "ASP") | .hostname' | while read line; do
        echo -e "for asp : \n $line"
        echo $line >> f1
    done

    ### jsp ###
    cat "result" | jq -r '. | select(.lang[0] == "JSP") | .hostname' | while read line; do
        echo -e "for jsp : \n $line"
        echo $line >> f1
    done

    ### python ###
    cat "result" | jq -r '. | select(.lang[0] == "PYTHON") | .hostname' | while read line; do
        echo -e "for python : \n $line"
        echo $line >> f1
    done



############## second filter #############
    cat "origin" | anew -d "f1" | while read line; do


    done

}

filter
