import Head from 'next/head';

export default function Home({ totalTips }) {
  return (
    <div style={{ fontFamily: 'Helvetica', margin: '0 auto', maxWidth: '550px', display: 'flex', flexDirection: 'column', alignItems: 'center', padding: '10px' }}>
      <Head>
        <title>neovim.tips</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <h1>neovim.tips API</h1>

      <div className="endpoint" style={{ marginBottom: '20px' }}>
        <h2>Get a Random Tip</h2>
        <strong>Endpoint:</strong> <code>/api/random</code><br />
        <strong>Method:</strong> <code>GET</code><br />
        <p>Returns a random tip from the collection of neovim tips.</p>
        <div className="curl-command" style={{ backgroundColor: 'black', color: 'white', padding: '2px 4px' }}>
          <code>curl -s https://www.neovim.tips/api/random</code>
        </div>
      </div>

      <div className="endpoint" id="specific-tip-example" style={{ marginBottom: '20px' }}>
        <h2>Get a Specific Tip by ID</h2>
        <strong>Endpoint:</strong> <code>{`/api/{1-${totalTips}}`}</code><br />
        <strong>Method:</strong> <code>GET</code><br />
        <p>Fetches a specific tip based on its unique identifier (ID).</p>
        <div className="curl-command" style={{ backgroundColor: 'black', color: 'white', padding: '2px 4px' }}>
          <code>curl -s https://www.neovim.tips/api/23</code>
        </div>
      </div>

      <footer>
        <p>Made with ❤️ in Nebraska</p>
      </footer>
    </div>
  );
}

export async function getServerSideProps(context) {
  try {
    const response = await fetch('https://www.neovim.tips/api/total');
    if (!response.ok) {
      throw new Error(`Failed to fetch: ${response.status}`);
    }
    const totalTips = await response.text();
    return { props: { totalTips } };
  } catch (error) {
    console.error('Error fetching total tips:', error);
    return { props: { totalTips: 'N/A' } };
  }
}
