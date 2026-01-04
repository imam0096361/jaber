package database

import (
	"log"
	"news-portal/src/model"
)

func InsertSampleData() {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM articles").Scan(&count)

	if count > 0 {
		return // Data already exists
	}

	articles := []model.Article{
		{
			Title:    "বাংলাদেশ অর্থনীতিতে নতুন মাইলফলক অর্জন",
			Content:  "বাংলাদেশের অর্থনীতি এই বছর নতুন উচ্চতায় পৌঁছেছে। জিডিপি বৃদ্ধির হার ৬.৫ শতাংশ অতিক্রম করেছে এবং রপ্তানি আয় রেকর্ড পরিমাণে বৃদ্ধি পেয়েছে।",
			Category: "জাতীয়",
			Author:   "রহিম আহমেদ",
			Featured: true,
		},
		{
			Title:    "ঢাকা ক্রিকেট ক্লাব আন্তর্জাতিক টুর্নামেন্টে চ্যাম্পিয়ন",
			Content:  "ঢাকা ক্রিকেট ক্লাব এশিয়া কাপ টুর্নামেন্টে বিজয়ী হয়েছে। দলটি ফাইনালে স্থানীয় প্রতিদ্বন্দ্বীদের পরাজিত করে চ্যাম্পিয়ন হওয়ার গৌরব অর্জন করেছে।",
			Category: "খেলাধুলা",
			Author:   "করিম খান",
			Featured: true,
		},
		{
			Title:    "নতুন প্রযুক্তি দিয়ে শিক্ষা ব্যবস্থা পরিবর্তন",
			Content:  "প্রযুক্তির কল্যাণে শিক্ষা ব্যবস্থা ক্রমাগত পরিবর্তিত হচ্ছে। কৃত্রিম বুদ্ধিমত্তা এবং মেশিন লার্নিং এখন শিক্ষার্থীদের ব্যক্তিগত শিক্ষণ পদ্ধতি তৈরি করতে সাহায্য করছে।",
			Category: "প্রযুক্তি",
			Author:   "আলি হোসেন",
			Featured: false,
		},
		{
			Title:    "স্বাস্থ্যকর জীবনযাত্রার জন্য নতুন গাইডলাইন",
			Content:  "বিশ্ব স্বাস্থ্য সংস্থা সম্প্রতি নতুন স্বাস্থ্যকর জীবনযাত্রার গাইডলাইন প্রকাশ করেছে। এতে দৈনিক ব্যায়াম, পুষ্টিকর খাবার এবং যথাযথ ঘুমের উপর জোর দেওয়া হয়েছে।",
			Category: "স্বাস্থ্য",
			Author:   "ডাক্তার সালমা",
			Featured: false,
		},
	}

	for _, article := range articles {
		query := `INSERT INTO articles (title, content, category, author, featured) VALUES (?, ?, ?, ?, ?)`
		_, err := DB.Exec(query, article.Title, article.Content, article.Category, article.Author, article.Featured)
		if err != nil {
			log.Println("Error inserting sample data:", err)
		}
	}
}
