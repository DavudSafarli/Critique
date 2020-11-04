package feedback_usecases

type FeedbackUsecases interface {
	FeedbackCreator
	FeedbackDetailsGetter
	FeedbackPaginator
}
