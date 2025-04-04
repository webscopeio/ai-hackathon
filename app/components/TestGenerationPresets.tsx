import React from 'react';
import { Card, CardContent } from "./ui/card";
import { Button } from '@/components/ui/button';
import { 
  LineChart, 
  ShieldCheck, 
  Workflow, 
  Lightbulb,
  Pencil
} from 'lucide-react';

export type PresetOption = {
  id: string;
  title: string;
  description: string;
  icon: React.ReactNode;
  prompt: string;
};

const presets: PresetOption[] = [
  {
    id: 'improve-coverage',
    title: 'Improve Test Coverage',
    description: 'Enhance overall test coverage by identifying gaps in the current test suite',
    icon: <ShieldCheck className="h-6 w-6 text-green-500" />,
    prompt: 'Analyze the codebase and identify areas with insufficient test coverage. Generate tests that improve the overall coverage, focusing on critical paths and edge cases.'
  },
  {
    id: 'sentry-issues',
    title: 'Cover Sentry Issues',
    description: 'Create tests for issues reported in Sentry to prevent regressions',
    icon: <LineChart className="h-6 w-6 text-red-500" />,
    prompt: 'Analyze Sentry error reports and create tests that would catch these issues. Focus on the most frequent errors and those affecting critical user flows.'
  },
  {
    id: 'user-flows',
    title: 'User Flow Coverage',
    description: 'Add tests for significant user journeys based on analytics',
    icon: <Workflow className="h-6 w-6 text-blue-500" />,
    prompt: 'Identify key user flows based on analytics data and create comprehensive tests for these journeys. Ensure that critical user paths are thoroughly tested.'
  },
  {
    id: 'custom',
    title: 'Custom Prompt',
    description: 'Create your own custom test generation prompt',
    icon: <Pencil className="h-6 w-6 text-gray-500" />,
    prompt: ''
  }
];

interface TestGenerationPresetsProps {
  onSelectPreset: (preset: PresetOption) => void;
  selectedPresetId: string | null;
}

export function TestGenerationPresets({ 
  onSelectPreset, 
  selectedPresetId 
}: TestGenerationPresetsProps) {
  return (
    <div className="space-y-4">
      <h3 className="text-lg font-medium flex items-center gap-2">
        <Lightbulb className="h-5 w-5 text-yellow-500" />
        Test Generation Presets
      </h3>
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {presets.map((preset) => (
          <Card 
            key={preset.id}
            className={`cursor-pointer transition-all hover:shadow-md ${
              selectedPresetId === preset.id 
                ? 'border-primary/40 ring-1 ring-primary/30' 
                : 'border-border'
            }`}
            onClick={() => onSelectPreset(preset)}
          >
            <CardContent className="p-4 flex flex-col h-full">
              <div className="flex items-start justify-between mb-2">
                <div className="p-2.5 rounded-full bg-muted">
                  {preset.icon}
                </div>
                {selectedPresetId === preset.id && (
                  <div className="bg-gray-100 p-1 rounded-full mt-1">
                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="text-primary">
                      <polyline points="20 6 9 17 4 12"></polyline>
                    </svg>
                  </div>
                )}
              </div>
              
              <h4 className="font-medium text-base mt-2">{preset.title}</h4>
              <p className="text-sm text-muted-foreground mt-2 flex-grow">
                {preset.description}
              </p>
              

            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
