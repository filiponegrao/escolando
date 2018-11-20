div=29
random=$(($RANDOM%$div))
random=$(($random+1))

bash get_text.sh ${random}