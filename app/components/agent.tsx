import { cn } from "@/lib/utils";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
} from "@/components/ui/card";
import { AgentTitle } from "@/components/agent-title";
import { Message } from "@/lib/types";
import ShinyText from "@/components/ShinyText";

type AgentProps = React.ComponentProps<typeof Card> & {
	active: boolean;
	messages: Message[];
	title: string;
	description: string;
};

export function Agent({
	className,
	active,
	messages,
	title,
	description,
	...props
}: AgentProps) {
	return (
		<Card
			className={cn("w-[380px] overflow-hidden", className)}
			active={active}
			{...props}
		>
			<CardHeader>
				<AgentTitle active={active}>{title}</AgentTitle>
				<CardDescription>
					{description + "\nClaude 3.5 Sonnet based Agent."}
				</CardDescription>
			</CardHeader>
			<CardContent className="flex flex-col gap-4">
				<div className="h-[300px] overflow-hidden relative">
					<div className="absolute bottom-0 w-full flex flex-col gap-1">
						{messages.map((message, index) => (
							<div
								key={index}
								className={cn(
									"mb-4 flex flex-col gap-1 items-start pb-4 last:mb-0 last:pb-0 transition-all duration-500 opacity-0",
									"animate-fade-in blur-sm",
								)}
							>
								<p className="text-sm font-medium leading-none capitalize">
									{message.title}
								</p>
								<p className="text-sm text-muted-foreground break-all">
									{message.description.length > 250
										? `${message.description.slice(0, 250)}...`
										: message.description}
								</p>
							</div>
						))}
					</div>
					<div
						className={cn(
							"absolute top-0 w-full h-[52px] bg-gradient-to-b from-background to-transparent",
							active && "from-blue-50/30",
						)}
					/>
				</div>
				<ShinyText
					className={cn(
						"justify-self-end",
						active ? "opacity-100" : "opacity-0",
					)}
					text="Generating..."
					speed={2}
				/>
			</CardContent>
		</Card>
	);
}
