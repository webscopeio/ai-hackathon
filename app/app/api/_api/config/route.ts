import { NextRequest, NextResponse } from "next/server";

// In a real application, this would be stored in a database
// For this hackathon project, we'll use a simple in-memory store
let configStore = {
  anthropicApiKey: "",
  sentryApiKey: "",
  techSpecification: "",
};

export async function GET() {
  return NextResponse.json(configStore);
}

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    
    // Update only the fields that are provided
    configStore = {
      ...configStore,
      ...body,
    };
    
    return NextResponse.json({ success: true, config: configStore });
  } catch (error) {
    return NextResponse.json(
      { error: "Failed to update configuration" },
      { status: 400 }
    );
  }
}
