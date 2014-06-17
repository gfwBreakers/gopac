var gulp = require('gulp'),
  concat = require('gulp-concat'),
  uglify = require('gulp-uglify');

gulp.task('default', function() {
  gulp.src(['./gogo.pac'])
    .pipe(concat('go.pac'))
    .pipe(uglify())
    .pipe(gulp.dest('.'));
});