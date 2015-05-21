var TagEditor = React.createClass({
    handleChange: function() {
        this.props.onUserInput(this.refs.tagsInput.getDOMNode().value);
    },
    render: function() {
        return (
                <div className="tag-editor" >
                <input
            type="text"
            name="tags"
            ref="tagsInput"
            value={this.props.tags}
            onChange={this.handleChange}
                >
                </input></div>
        );
    }
})

var RandomButton = React.createClass({
    render: function() {
        return (
                <button onClick={this.props.onClick}>Random Article!</button>
        );
    }
})

var SearchArticlesButton = React.createClass({
    render: function() {
        return (
                <button onClick={this.props.onClick}>List Articles!</button>
        );
    }
})

var Article = React.createClass({
    render: function() {
        var article = this.props.article;
        if (!article.title) {
            return (<div className="garticle"></div>);
        }
        var tagNodes = (article.tags?article.tags:[]).map(function(tag){
            return(
                    <span key={tag}>{tag} | </span>
            );
        });

        return (
                <article className="garticle">
                <h1> {article.title} </h1>
                <div className="garticle-metadata">
                <span className="garticle-author"> {article.author} </span> |
                <span className="garticle-blog"> {article.blog} </span> |
                <span className="article-date"> <a href={article.url}> {article.date_published} </a> </span>
                </div>
                <div className="garticle-content" dangerouslySetInnerHTML={{__html: article.content}} />
                <div className="garticle-tags"> {tagNodes}</div>
                </article>
        );
    }
})


var ArticleList = React.createClass({
    render: function() {
        var articles = this.props.articles;
        if (articles.length == 0) {
            return (<div className="garticle-list"></div>);
        }
        var articleNodes = articles.map(function(article){
            var tags = article.tags || []
            var read = tags.indexOf("read") > -1?"garticle-read":""
            var class_name = "garticle-info " + read
            return(
                    <article className={class_name} key={article.url}>
                    <div className="garticle-metadata">
                    <div className="garticle-blog"> {article.blog} </div>
                    <div className="garticle-title"> {article.title} </div>
                    </div>
                    </article>
            );
        });

        var height = $(document).height() * 0.95;
        var style =  {height: height + "px"}
        return (
                <div className="garticle-list" style={style}>
                {articleNodes}
                </div>
        );
    }
})


var GooverApp = React.createClass({
    getInitialState: function(){
        return (
            {
                article: {},
                tags: "!read",
                articleList: []
            }
        );
    },
    updateTags: function(tags){
        this.setState({tags: tags});
    },
    fetchArticle: function(){
        var self = this;
        var fetch = $.get("/random?tag=" + this.state.tags)
            .done(function (article) {
                self.setState({article: article["entry"], articleList: []});
            })
            .error(function (data, response) {
                self.setState({article: {}, articles:[]});
            })

    },
    listArticles: function(){
        var self = this;
        var fetch = $.get("/view?tag=" + this.state.tags)
            .done(function (data) {
                console.log(data);
                if (!data) {
                    data = []
                }
                self.setState({articleList: data, article: {}});
            })
            .error(function (data, response) {
                self.setState({article: {}, articles:[]});
            })
    },
    render: function() {
        return (
                <div className="gooverapp">
                <TagEditor tags={this.state.tags} onUserInput={this.updateTags}/>
                <RandomButton onClick={this.fetchArticle} />
                <SearchArticlesButton onClick={this.listArticles}/>
                <Article article={this.state.article} />
                <ArticleList articles={this.state.articleList} />
                </div>
        );
    }
})

// Implement a TagEditor that can be used both in an article and a list of
// articles.  Currently, /edit/ end-point just allows editing the tags of one
// article. We could have an end point that supports editing tags for multiple
// articles (may be a different end-point.)

React.render(
        <GooverApp />,
  document.getElementById('content')
);
