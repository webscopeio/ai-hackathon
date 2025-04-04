"use client";

import React, { useState } from "react";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Sparkles, ArrowRight } from "lucide-react";
import { toast } from "sonner";
import { TestGenerationPresets, PresetOption } from "./TestGenerationPresets";
import { Card, CardContent } from "./ui/card";

interface TestGeneratorProps {
  techSpecification: string;
  productSpecification: string;
  onGenerateTests: (prompt: string) => Promise<void>;
}

export function TestGenerator({
  techSpecification,
  productSpecification,
  onGenerateTests
}: TestGeneratorProps) {
  const [selectedPreset, setSelectedPreset] = useState<PresetOption | null>(null);
  const [customPrompt, setCustomPrompt] = useState("");
  const [isGenerating, setIsGenerating] = useState(false);
  
  // Handle preset selection
  const handlePresetSelect = (preset: PresetOption) => {
    setSelectedPreset(preset);
    
    // If it's the custom preset, clear the prompt to let user enter their own
    if (preset.id === 'custom') {
      setCustomPrompt("");
    } else {
      // Otherwise use the preset's prompt
      setCustomPrompt(preset.prompt);
    }
  };
  
  // Handle custom prompt changes
  const handlePromptChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setCustomPrompt(e.target.value);
  };
  
  // Handle test generation
  const handleGenerateTests = async () => {
    if (!selectedPreset) {
      toast.error("Please select a test generation preset or custom prompt");
      return;
    }
    
    if (selectedPreset.id === 'custom' && !customPrompt.trim()) {
      toast.error("Please enter a custom prompt for test generation");
      return;
    }
    
    try {
      setIsGenerating(true);
      
      // Use the final prompt (either preset or custom)
      const finalPrompt = customPrompt;
      
      await onGenerateTests(finalPrompt);
      toast.success("Tests generated successfully");
    } catch (error) {
      console.error("Error generating tests:", error);
      toast.error("Failed to generate tests");
    } finally {
      setIsGenerating(false);
    }
  };
  
  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2">
        <Sparkles className="h-5 w-5 text-primary" />
        <h2 className="text-xl font-semibold">Generate Tests</h2>
      </div>
      
      <p className="text-muted-foreground">
        Select a preset below to generate tests for your application. Each preset focuses on different aspects of testing to help improve your test coverage and quality.
      </p>
      
      <TestGenerationPresets 
        onSelectPreset={handlePresetSelect}
        selectedPresetId={selectedPreset?.id || null}
      />
      
      {selectedPreset?.id === 'custom' && (
        <div className="space-y-2">
          <div className="flex justify-between items-center">
            <h3 className="text-sm font-medium">Custom Test Generation Prompt</h3>
            <Button 
              variant="ghost" 
              size="sm"
              onClick={() => setSelectedPreset(null)}
            >
              Clear Selection
            </Button>
          </div>
          
          <Textarea
            placeholder="Enter your custom test generation instructions..."
            className="min-h-24"
            value={customPrompt}
            onChange={handlePromptChange}
            disabled={isGenerating}
          />
          
          <p className="text-xs text-muted-foreground">
            Provide specific instructions for test generation
          </p>
        </div>
      )}
      

      
      <div className="flex justify-end">
        <Button 
          onClick={handleGenerateTests} 
          disabled={isGenerating || !selectedPreset}
          className="gap-2"
        >
          {isGenerating ? "Generating..." : "Generate Tests"}
          {!isGenerating && <ArrowRight className="h-4 w-4" />}
        </Button>
      </div>
    </div>
  );
}
