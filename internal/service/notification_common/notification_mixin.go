// notification mixin
package notificationcommon

import (
	"context"
	"encoding/base64"
	"regexp"

	"github.com/apache/incubator-answer/internal/service/mixinbot"
	mixinbotlang "github.com/apache/incubator-answer/internal/service/mixinbot/lang"
	"github.com/apache/incubator-answer/plugin"
	"github.com/goccy/go-json"

	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/segmentfault/pacman/log"
)

func (ns *NotificationCommon) sendMixinNotification(notificationMsgDescription string, msg *plugin.NotificationMessage) error {
	notificationMsgTitle, notificationMsgContent := mixinbot.PrefixTitle+msg.QuestionTitle, msg.Content

	titleRunes := []rune(notificationMsgTitle)
	if len(titleRunes) > mixinbot.MaxCardTitleLength {
		notificationMsgTitle = string(titleRunes[:mixinbot.MaxCardTitleLength-len([]rune(mixinbot.Ellipsis))]) + mixinbot.Ellipsis
	}

	re := regexp.MustCompile(`<.*?>`)
	notificationMsgContent = re.ReplaceAllString(notificationMsgContent, "")

	contentRunes := []rune(notificationMsgContent)
	if len(contentRunes) > mixinbot.MaxCardContentLength {
		notificationMsgContent = string(contentRunes[:mixinbot.MaxCardContentLength-len([]rune(mixinbot.Ellipsis))]) + mixinbot.Ellipsis
	}

	lastContent := "\n\n\n" + notificationMsgDescription
	totalContentRunes := []rune(notificationMsgContent + lastContent)
	if len(totalContentRunes) <= mixinbot.MaxCardContentLength {
		notificationMsgContent += lastContent
	} else {
		availableLength := mixinbot.MaxCardContentLength - len([]rune(lastContent)) - len([]rune(mixinbot.Ellipsis))
		notificationMsgContent = string(contentRunes[:availableLength]) + mixinbot.Ellipsis + lastContent
	}

	card := &mixin.AppCardMessage{
		AppID:       ns.mixinbotService.Config.ClientID,
		Title:       notificationMsgTitle,
		Description: notificationMsgContent,
		Shareable:   true,
	}
	ns.fillCardAction(card, msg)

	cardBytes, err := json.Marshal(card)
	if err != nil {
		return err
	}

	cardBase64code := base64.StdEncoding.EncodeToString(cardBytes)
	messageRequest := &mixin.MessageRequest{
		ConversationID: mixin.UniqueConversationID(ns.mixinbotService.Config.ClientID, msg.ReceiverExternalID),
		RecipientID:    msg.ReceiverExternalID,
		MessageID:      mixin.RandomTraceID(),
		Category:       mixin.MessageCategoryAppCard,
		Data:           cardBase64code,
	}

	return ns.mixinbotService.SendMessage(context.Background(), messageRequest)
}

func (ns *NotificationCommon) fillCardAction(card *mixin.AppCardMessage, msg *plugin.NotificationMessage) {
	btnMsg := mixin.AppButtonMessage{
		Label:  ns.langPicker.Pick(mixinbotlang.GetLanguage(msg.ReceiverLang)).TranslateGetCardInfo(),
		Action: msg.QuestionUrl,
		Color:  mixinbot.RandomCardColor(),
	}

	switch msg.Type {
	case plugin.NotificationUpdateQuestion, plugin.NotificationInvitedYouToAnswer, plugin.NotificationNewQuestion, plugin.NotificationNewQuestionFollowedTag:
		btnMsg.Action = msg.QuestionUrl
	case plugin.NotificationAnswerTheQuestion, plugin.NotificationUpdateAnswer, plugin.NotificationAcceptAnswer:
		btnMsg.Action = msg.AnswerUrl
	case plugin.NotificationCommentQuestion, plugin.NotificationCommentAnswer, plugin.NotificationReplyToYou, plugin.NotificationMentionYou:
		btnMsg.Action = msg.CommentUrl
	default:
		log.Debugf("this type of notification will be drop, the type is %s", msg.Type)
	}
	card.Actions = []mixin.AppButtonMessage{btnMsg}
}
