package handlers

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"start-feishubot/initialization"
)

type CardKind string
type CardChatType string

var (
	ClearCardKind = CardKind("clear")
)

var (
	GroupChatType = CardChatType("group")
	UserChatType  = CardChatType("user")
)

type CardMsg struct {
	kind     CardKind
	chatType CardChatType
	value    interface{}
}

func sendCard(ctx context.Context,
	chatId *string,
	cardContent string,
) error {
	client := initialization.GetLarkClient()
	resp, err := client.Im.Message.Create(ctx, larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeChatId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			ReceiveId(*chatId).
			Content(cardContent).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return err
	}
	return nil
}

func newSendCard(elements ...larkcard.MessageCardElement) (string,
	error) {
	config := larkcard.NewMessageCardConfig().
		WideScreenMode(false).
		EnableForward(true).
		UpdateMulti(false).
		Build()
	// header
	title := "🤖️机器人提醒"
	header := larkcard.NewMessageCardHeader().
		Template(larkcard.TemplateBlue).
		Title(larkcard.NewMessageCardPlainText().
			Content(title).
			Build()).
		Build()

	var aElementPool []larkcard.MessageCardElement
	for _, element := range elements {
		aElementPool = append(aElementPool, element)
	}
	// 卡片消息体
	cardContent, err := larkcard.NewMessageCard().
		Config(config).
		Header(header).
		Elements(
			aElementPool,
		).
		String()
	return cardContent, err
}

func withNote(note string) larkcard.MessageCardElement {
	noteElement := larkcard.NewMessageCardNote().
		Elements([]larkcard.MessageCardNoteElement{larkcard.NewMessageCardPlainText().
			Content("请注意，这将开始一个全新的对话，您将无法利用之前的对话历史信息").
			Build()}).
		Build()
	return noteElement
}

func withMainMsg(msg string) larkcard.MessageCardElement {
	msg, i := processMessage(msg)
	if i != nil {
		return nil
	}
	mainElement := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{larkcard.NewMessageCardField().
			Text(larkcard.NewMessageCardLarkMd().
				Content(msg).
				Build()).
			IsShort(true).
			Build()}).
		Build()
	return mainElement
}

func withDoubleCheckBtn() larkcard.MessageCardElement {
	actions := larkcard.NewMessageCardAction().
		Actions([]larkcard.MessageCardActionElement{
			larkcard.NewMessageCardEmbedButton().
				Type(larkcard.MessageCardButtonTypeDanger).
				Value(map[string]interface{}{
					"clearCheck": 1, // 1 代表确认清除
					"kind":       ClearCardKind,
					"chatType":   UserChatType,
				}).
				Text(larkcard.NewMessageCardPlainText().
					Content("确认清除").
					Build()),
			larkcard.NewMessageCardEmbedButton().
				Type(larkcard.MessageCardButtonTypePrimary).
				Value(map[string]interface{}{
					"clearCheck": 0, // 0 代表取消清除
					"kind":       ClearCardKind,
					"chatType":   UserChatType,
				}).
				Text(larkcard.NewMessageCardPlainText().
					Content("我再想想").
					Build()),
		}).Layout(larkcard.MessageCardActionLayoutBisected.Ptr()).
		Build()
	return actions
}
func replyMsg(ctx context.Context, msg string, msgId *string) error {
	fmt.Println("sendMsg", msg, msgId)
	msg, i := processMessage(msg)
	if i != nil {
		return i
	}
	client := initialization.GetLarkClient()
	content := larkim.NewTextMsgBuilder().
		Text(msg).
		Build()

	resp, err := client.Im.Message.Reply(ctx, larkim.NewReplyMessageReqBuilder().
		MessageId(*msgId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			Uuid(uuid.New().String()).
			Content(content).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return err
	}
	return nil
}
func sendMsg(ctx context.Context, msg string, chatId *string) error {
	//fmt.Println("sendMsg", msg, chatId)
	msg, i := processMessage(msg)
	if i != nil {
		return i
	}
	client := initialization.GetLarkClient()
	content := larkim.NewTextMsgBuilder().
		Text(msg).
		Build()

	//fmt.Println("content", content)

	resp, err := client.Im.Message.Create(ctx, larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeChatId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			ReceiveId(*chatId).
			Content(content).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return err
	}
	return nil
}

func sendClearCacheCheckCard(ctx context.Context, chatId *string) {
	newCard, _ := newSendCard(
		withMainMsg("您确定要清除对话上下文吗？"),
		withNote("请注意，这将开始一个全新的对话，您将无法利用之前的对话历史信息"),
		withDoubleCheckBtn())
	sendCard(
		ctx,
		chatId,
		newCard,
	)
}
