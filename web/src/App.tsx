import { useState, useEffect } from "react";

import "./App.css";

type Tweet = {
  id: number;
  text: string;
  favourites: number;
  retweets: number;
  date: string;
  device: string;
  isRetweet: boolean;
  isDeleted: boolean;
  isFlagged: boolean;
  isReal: boolean;
};

type MaybeTweet = Tweet | null;

function TweetCard({ tweet }: { tweet: MaybeTweet }) {
  if (!tweet) {
    return <div>not yet...</div>;
  }
  return (
    <div className="card">
      <p>{tweet.text}</p>
      <p>{tweet.date}</p>
      <p>{tweet.device}</p>
      <p>{tweet.favourites} likes</p>
      <p>{tweet.retweets} retweets</p>

      <p class="spoiler">Real or fake: {tweet.isReal ? "Real" : "Fake"}</p>
    </div>
  );
}

type tObj = {
  id: number;
  text: string;
  favorites: number;
  retweets: number;
  date: string;
  device: string;
  isRetweet: boolean;
  isDeleted: boolean;
  isFlagged: boolean;
  isReal: boolean;
};

function parseToTweet(data: tObj): Tweet {
  return {
    id: data.id,
    text: data.text,
    favourites: data.favorites,
    retweets: data.retweets,
    date: data.date,
    device: data.device,
    isRetweet: data.isRetweet,
    isDeleted: data.isDeleted,
    isFlagged: data.isFlagged,
    isReal: data.isReal,
  };
}

function App() {
  const [tweet, setTweet] = useState<Tweet | null>(null);

  useEffect(() => {
    fetch("http://localhost:1776/v1/trump")
      .then((response) => response.json())
      .then((data) => {
        const first = data[0];
        console.log(first);
        const t = parseToTweet(first);
        setTweet(t);
      });
  }, []);

  return (
    <>
      <h1>Trump or Markov?</h1>
      <TweetCard tweet={tweet} />
    </>
  );
}

export default App;
