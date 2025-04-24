"use client";

import { useState, useEffect, useRef } from "react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";
import { Message } from "@/lib/types";
import { Agent } from "@/components/agent";
import { Tool } from "@/components/tool";
import { TestFiles } from "./test-files";
export default function WebSocketDemo() {
  const [messages, setMessages] = useState<string[]>([]);
  const [scenarioCount, setScenarioCount] = useState<number>(0);
  const [fileNames, setFileNames] = useState<string[]>([]);
  const [scenario, setScenario] = useState<string>("");
  const [inputValue, setInputValue] = useState(
    "https://ai-hackathon-demo-delta.vercel.app/"
  );
  const [isConnected, setIsConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);
  const scrollAreaRef = useRef<HTMLDivElement>(null);

  const addMessage = (message: string) => {
    setMessages((prev) => [...prev, message]);
  };

  const handleClose = () => {
    if (!wsRef.current) return;
    wsRef.current.close();
  };

  const handleSend = () => {
    if (!wsRef.current) {
      try {
        console.log("Attempting to connect to WebSocket...");
        wsRef.current = new WebSocket("ws://localhost:8080/ws");

        wsRef.current.onopen = () => {
          console.log("WebSocket connection opened successfully");
          setIsConnected(true);
          addMessage("OPEN");
          addMessage(`SEND '${inputValue}'`);
          if (wsRef.current) {
            wsRef.current.send(inputValue);
          }
        };

        wsRef.current.onclose = (event) => {
          console.log("WebSocket connection closed:", event.code, event.reason);
          setIsConnected(false);
          addMessage(
            `CLOSE '${event.code}${
              event.reason ? `, Reason: ${event.reason}` : ""
            }'`
          );
          wsRef.current = null;
        };

        wsRef.current.onmessage = (evt) => {
          console.log(evt.data);
          addMessage(evt.data);
          const id = evt.data.split(" ")[0];
          if (id === "SCENARIOS") {
            const numScenarios = parseInt(evt.data.split(" ")[1]);
            setScenarioCount(numScenarios);
          }
          if (id === "FILENAME") {
            const fileName = evt.data.split(" ")[1];
            setFileNames((prev) => [...prev, fileName]);
          }
          if (id === "SCENARIO") {
            const scenario = evt.data.substring(evt.data.indexOf(" ") + 1);
            setScenario(scenario);
          }
        };

        wsRef.current.onerror = (error) => {
          console.error("WebSocket error:", error);
          addMessage(
            `ERROR 'Connection failed - check browser console for details'`
          );
        };
      } catch (error: any) {
        console.error("WebSocket connection error:", error);
        addMessage(
          `ERROR 'Failed to connect to WebSocket server - ${error.message}'`
        );
      }
      return;
    }
    addMessage(`SEND: ${inputValue}`);
    wsRef.current.send(inputValue);
  };

  useEffect(() => {
    if (scrollAreaRef.current) {
      scrollAreaRef.current.scrollTop = scrollAreaRef.current.scrollHeight;
    }
  }, [messages]);

  return (
    <div className="container mx-auto p-6">
      <Card className="w-4xl mx-auto mb-12">
        <CardHeader className="pb-0">
          <CardTitle>End-to-end Test Generation</CardTitle>
          <CardDescription>
            Generate end-to-end tests that are passing with just a prompt and
            the website URL.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col gap-4">
            <div className="space-y-4">
              <div className="flex gap-2"></div>
              <div className="flex gap-2">
                <Textarea
                  value={inputValue}
                  onChange={(e) => setInputValue(e.target.value)}
                  placeholder="Write a the website's URL, yeah, that's it."
                />
                <div className="flex flex-col gap-2">
                  {" "}
                  <Button onClick={handleSend} variant="secondary">
                    Send
                  </Button>
                  <Button
                    onClick={handleClose}
                    disabled={!isConnected}
                    variant="destructive"
                  >
                    Close
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
      <div className="flex flex-row gap-8 justify-center">
        <Agent
          title="Analyzer"
          description="Analyzes the website and generates scenarios."
          active={
            isActive(messages, "ANALYZER") || isActive(messages, "TOOLCALL")
          }
          messages={getMessages(messages, ["ANALYZER", "TOOLCALL"])}
        />
        <Agent
          title="Generator"
          description="Generates end-to-end test from scenarios."
          active={isActive(messages, "GENERATOR")}
          messages={getMessages(messages, ["GENERATOR"])}
        />
        <Agent
          title="Evaluator"
          description="Evaluates the generated end-to-end test."
          active={isActive(messages, "EVALUATOR")}
          messages={getMessages(messages, ["EVALUATOR"])}
        />
      </div>
      <div className="flex flex-row gap-8 justify-center pt-16">
        <Tool
          title="get_sitemap_tool"
          description="Gets the content of the website."
          active={isActive(messages, "GET_SITEMAP_TOOL")}
          messages={getMessages(messages, ["GET_SITEMAP_TOOL"])}
        />
        <Tool
          title="get_content_tool"
          description="Gets the content of passed urls."
          active={isActive(messages, "GET_CONTENT_TOOL")}
          messages={getMessages(messages, ["GET_CONTENT_TOOL"])}
        />
        <div className="w-[470px]"></div>
        <Tool
          title="run_test_tool"
          description="Runs the generated end-to-end test."
          active={isActive(messages, "RUN_TEST_TOOL")}
          messages={getMessages(messages, ["RUN_TEST_TOOL"])}
        />
      </div>
      <TestFiles
        className="pt-24"
        count={scenarioCount}
        names={fileNames}
        scenario={scenario}
      />
    </div>
  );
}

const isActive = (messages: string[], id: string): boolean => {
  return (
    messages.length > 0 && messages[messages.length - 1].split(" ")[0] === id
  );
};

const getMessages = (messages: string[], ids: string[]): Message[] => {
  const agentName = ids[0];
  return messages
    .filter((message) => ids.includes(message.split(" ")[0]))
    .map((message) => {
      const id = message.split(" ")[0];
      const description = message.split(" ").slice(1).join(" ");
      return { title: id === agentName ? "Agent" : "Toolcall", description };
    });
};
