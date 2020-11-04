"use strict";
$(document).ready(function() {


    /*------------ Start site menu  ------------*/

    // Start sticky header
    $(window).on('scroll', function() {
        if ($(window).scrollTop() >= 150) {
            $('#sticky-header').addClass('sticky-menu');
        } else {
            $('#sticky-header').removeClass('sticky-menu');
        }
    });

    // slicknav
    $('ul#navigation').slicknav({
        prependTo: ".responsive-menu-wrap"
    });


    // cart script
    $('.add').click(function() {
        if ($(this).prev().val() < 10000000) {
            $(this).prev().val(+$(this).prev().val() + 1);
        }
    });
    $('.sub').click(function() {
        if ($(this).next().val() > 1) {
            if ($(this).next().val() > 1) $(this).next().val(+$(this).next().val() - 1);
        }
    });

    $('.search-close').on('click', function() {
        $('.gold-search-form').addClass('gold-search-hide');
    });

    $('.gold-search').on('click', function() {
        $('.gold-search-form').removeClass('gold-search-hide');
    });



    /* Product light */

    $('.gold-product-carousel').owlCarousel({
        loop: true,
        margin: 10,
        nav: false,
        items: 1,
    });





});