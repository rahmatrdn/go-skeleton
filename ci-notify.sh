TIME="10"
URL="https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/sendMessage"
TEXT="Deploy status: $1 $2%0A%0AProject:+$PROJECT_NAME+-+$CI_PROJECT_NAME%0AURL:+$CI_PROJECT_URL/pipelines/$CI_PIPELINE_ID/%0ABranch:+$CI_COMMIT_REF_SLUG%0AAuthor:+$CI_COMMIT_AUTHOR%0APipelineID:+$CI_PIPELINE_ID"
#TEXT="Deploy status: $1 %0A%0AProject:+$PROJECT_NAME+-+$CI_PROJECT_NAME%0AURL:+$CI_PROJECT_URL/pipelines/$CI_PIPELINE_ID/%0ABranch:+$CI_COMMIT_REF_SLUG%0ASource Branch:+$3%0AReviewer:+$CI_COMMIT_AUTHOR%0AAssignee:+$2%0APipelineID:+$CI_PIPELINE_ID"
curl -s --max-time $TIME -d "chat_id=$TELEGRAM_USER_ID&disable_web_page_preview=1&text=$TEXT" $URL > /dev/null