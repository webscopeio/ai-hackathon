interface FeedbackProps {
  feedback: string;
  className?: string;
}

export default function Feedback({ feedback, className }: FeedbackProps) {
  if (!feedback) {
    return null;
  }

  return (
    <div className={className}>
      <h2 className="text-lg font-semibold mb-2">Current feedback:</h2>
      <p>{feedback}</p>
    </div>
  );
}
