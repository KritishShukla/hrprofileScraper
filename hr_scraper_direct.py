import os
import re
from tavily import TavilyClient
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

def main():
    # Check if API key is set
    api_key = os.getenv("TAVILY_API_KEY")
    if not api_key:
        print("Error: TAVILY_API_KEY not found in environment variables.")
        print("Please set your Tavily API key in a .env file or as an environment variable.")
        return
    
    # Define your target company here
    company_name = input("Enter the company name (or press Enter for 'Microsoft'): ").strip()
    if not company_name:
        company_name = "Microsoft"
    
    print(f"\n🔍 Searching for HR profiles at {company_name}...")
    
    try:
        # Initialize the Tavily client
        tavily = TavilyClient(api_key=api_key)
        
        # Define the search query with the company name
        # This query looks for common HR-related titles at the specified company
        query = f"('Human Resources' OR 'Talent Acquisition' OR 'Recruiter') '{company_name}' site:linkedin.com/in/"
        
        print(f"Executing search with query: {query}\n")
        
        # Execute the search
        response = tavily.search(query=query, max_results=20)
        
        # Extract and print the LinkedIn URLs
        print(f"Found LinkedIn Profile URLs for HR at {company_name}:")
        linkedin_urls = []
        
        for result in response.get('results', []):
            url = result.get('url')
            # Ensure the URL is a valid LinkedIn profile URL
            if url and "linkedin.com/in/" in url:
                # Clean the URL to remove tracking parameters
                clean_url = url.split('?')[0]
                if clean_url not in linkedin_urls:
                     linkedin_urls.append(clean_url)
        
        # Print the final list of unique URLs
        if linkedin_urls:
            print(f"\n✅ Found {len(linkedin_urls)} unique LinkedIn profiles:")
            for i, url in enumerate(linkedin_urls, 1):
                print(f"{i}. {url}")
        else:
            print("No LinkedIn profiles found for the specified search.")
            print("This might be due to:")
            print("- No HR professionals found at this company on LinkedIn")
            print("- LinkedIn blocking search results")
            print("- Company name not matching LinkedIn profiles")
            
    except Exception as e:
        print(f"Error occurred during search: {str(e)}")
        print("Please check your API key and internet connection.")

if __name__ == "__main__":
    main()
