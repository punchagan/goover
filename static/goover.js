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

var Article = React.createClass({
    render: function() {
        var article = this.props.article;
        if (!article.title) {
            return (<div className="garticle">no article!</div>);
        }
        var tagNodes = article.tags.map(function(tag){
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


var GooverApp = React.createClass({
    getInitialState: function(){
        return (
            {
                article: {},
                tags: "!read"
            }
        );
    },
    updateTags: function(tags){
        this.setState({tags: tags});
    },
    fetchArticle: function(){
        var self = this;
        var fetch = $.get("/random?tags=" + this.state.tags)
            .done(function (article) {
                console.log(article);
                self.setState({article: article});
            })
            .error(function (data, response) {
                self.setState({article: {title: "No article found", tags:self.state.tags.split(',')}});
            })

    },
    render: function() {
        return (
                <div className="gooverapp">
                <TagEditor tags={this.state.tags} onUserInput={this.updateTags}/>
                <RandomButton onClick={this.fetchArticle} />
                <Article article={this.state.article} />
                </div>
        );
    }
})


React.render(
        <GooverApp />,
  document.getElementById('content')
);
